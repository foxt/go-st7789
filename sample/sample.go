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
	showWork := make(chan image.Image, 3)
	waitGroup := sync.WaitGroup{}
	displayCtx, cancelFunc := context.WithCancel(ctx)
	waitGroup.Add(1)
	go func() {
		defer func() {
			waitGroup.Done()
		}()
		for {
			select {
			case <-displayCtx.Done():
				return
			case img := <-showWork:
				canvas.DrawImage(img)
				canvas.Flush()
			}
		}
	}()
	defer func() {
		cancelFunc()
		waitGroup.Wait()
	}()
	for {
		for i, img := range all.Image {
			if ctx.Err() != nil {
				return
			}
			showWork <- img
			time.Sleep(time.Duration(all.Delay[i]) * time.Millisecond * 10)
		}
		if all.LoopCount < 0 {
			break
		}
		if all.LoopCount != 0 {
			all.LoopCount -= 1
		}
	}
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
	device := ST7789.NewST7789(
		rpio.Spi0,
		rpio.Pin(25),
		rpio.Pin(27),
		rpio.Pin(24),
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
