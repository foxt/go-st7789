# ST7789

使用Goland实现的操作ST7789,适用于无CS引脚的240x204 LCD显示屏。 目前仅在Raspberry zero 2w上测试通过。

本库根据Python版 https://github.com/solinnovay/Python_ST7789 移植而来，并在此基础上实现了RGBA转RGB565,支持透明图层。
# 安装
```shell
go get github.com/manx98/go-st7789
```
# 使用示例
```go
package main

import (
	"context"
	ST7789 "github.com/manx98/go-st7789"
	"github.com/stianeikeland/go-rpio/v4"
	"image"
	"image/gif"
	"log"
	"os"
	"sync"
	"time"
)

// displayGIF
//
//	@Description: 显示GIF图片
//	@param ctx
//	@param canvas 画布
//	@param filePath GIF路径
func displayGIF(ctx context.Context, canvas *ST7789.Canvas, filePath string) {
	f, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to open: %v", err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatalf("failed to close: %v", err)
		}
	}()
	all, err := gif.DecodeAll(f)
	if err != nil {
		log.Fatalf("failed to decode: %v", err)
	}
	showWork := make(chan image.Image, 1)
	waitGroup := sync.WaitGroup{}
	displayCtx, cancelFunc := context.WithCancel(ctx)
	waitGroup.Add(1)
	var totalTime int64
	var total int64
	go func() {
		defer func() {
			waitGroup.Done()
			cancelFunc()
		}()
		for {
			select {
			case <-displayCtx.Done():
				return
			case img := <-showWork:
				start := time.Now()
				canvas.DrawImage(img)
				canvas.Flush()
				totalTime += time.Now().Sub(start).Milliseconds()
				total++
			}
		}
	}()
	defer func() {
		cancelFunc()
		waitGroup.Wait()
		log.Printf("平均速度：%dms/fps\n", totalTime/total)
	}()
	for {
		for i, img := range all.Image {
			select {
			case <-displayCtx.Done():
				return
			case showWork <- img:
				time.Sleep(time.Duration(all.Delay[i]) * time.Millisecond * 10)
			}
		}
		if all.LoopCount < 0 {
			break
		}
		if all.LoopCount != 0 {
			all.LoopCount -= 1
		}
	}
}

type MyPin struct {
	rpio.Pin
}

func (m *MyPin) SetOutput() {
	m.Mode(rpio.Output)
}

type MySpi struct {
}

func (m *MySpi) SpiSpeed(speed uint32) {
	rpio.SpiSpeed(int(speed))
}

func (m *MySpi) SetSpiMode3() {
	rpio.SpiMode(1, 1)
}

func (m *MySpi) SpiTransmit(data []byte) {
	rpio.SpiExchange(data)
}

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio: %v", err)
	}
	defer func() {
		if err := rpio.Close(); err != nil {
			log.Fatalf("failed to close gpio: %v", err)
		}
	}()
	err := rpio.SpiBegin(rpio.Spi0)
	if err != nil {
		log.Fatalf("failed to begin gpio: %v", err)
	}
	device := ST7789.NewST7789(
		&MySpi{},
		&MyPin{rpio.Pin(25)},
		&MyPin{rpio.Pin(27)},
		&MyPin{rpio.Pin(24)},
		ST7789.ST7789_TFTWIDTH,
		ST7789.ST7789_TFTHEIGHT,
	)
	canvas := device.GetFullScreenCanvas()
	timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	displayGIF(timeout, canvas, "./sample/TeaTime.gif")
	canvas.Clear()
	canvas.Flush()
}
```
# 感谢
1. Python原始实现 https://github.com/solinnovay/Python_ST7789
2. GPIO库 https://github.com/stianeikeland/go-rpio/
