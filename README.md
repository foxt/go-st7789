# ST7789

> [!NOTE]
> This readme, and some code comments was machine translated from the original Chinese - https://github.com/manx98/go-st7789
> Some inaccuracies may be present.

ST7789 library implemented using Golang, suitable for 240x204 LCD display without CS pin. Currently only tested on Raspberry zero 2w.


This library is transplanted from the Python version https://github.com/solinnovay/Python_ST7789, and based on this, it implements RGBA conversion to RGB565 and supports transparent layers.

# Installation
```shell
go get github.com/manx98/go-st7789
```
# Usage example
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
//	@Description: Display GIF image
//	@param ctx
//	@param canvas 
//	@param filePath path to the GIF file
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
		log.Printf("Average speed：%dms/frame\n", totalTime/total)
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
	rpio.SpiTransmit(data...)
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
		ST7789.Screen240X240,
	)
	canvas := device.GetFullScreenCanvas()
	timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	displayGIF(timeout, canvas, "./sample/TeaTime.gif")
	canvas.Clear()
	canvas.Flush()
}
```
# Acknowledgements
1. Python original implementation https://github.com/solinnovay/Python_ST7789
2. GPIO library https://github.com/stianeikeland/go-rpio/
