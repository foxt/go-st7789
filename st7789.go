package ST7789

import (
	"image"
	"image/color"
	"sync"
	"time"
)

const (
	SPI_CLOCK_HZ = 40000000 // 40 MHz

	// Constants for interacting width display registers.
	ST7789_TFTWIDTH  = 240
	ST7789_TFTHEIGHT = 240

	ST7789_NOP       = 0x00
	ST7789_SWRESET   = 0x01
	ST7789_RDDID     = 0x04
	ST7789_RDDST     = 0x09
	ST7789_RDDPM     = 0x0A
	ST7789_RDDMADCTL = 0x0B
	ST7789_RDDCOLMOD = 0x0C
	ST7789_RDDIM     = 0x0D
	ST7789_RDDSM     = 0x0E
	ST7789_RDDSDR    = 0x0F

	ST7789_SLPIN  = 0x10
	ST7789_SLPOUT = 0x11
	ST7789_PTLON  = 0x12
	ST7789_NORON  = 0x13

	ST7789_INVOFF  = 0x20
	ST7789_INVON   = 0x21
	ST7789_GAMSET  = 0x26
	ST7789_DISPOFF = 0x28
	ST7789_DISPON  = 0x29
	ST7789_CASET   = 0x2A
	ST7789_RASET   = 0x2B
	ST7789_RAMWR   = 0x2C
	ST7789_RAMRD   = 0x2E

	ST7789_PTLAR    = 0x30
	ST7789_VSCRDEF  = 0x33
	ST7789_TEOFF    = 0x34
	ST7789_TEON     = 0x35
	ST7789_MADCTL   = 0x36
	ST7789_VSCRSADD = 0x37
	ST7789_IDMOFF   = 0x38
	ST7789_IDMON    = 0x39
	ST7789_COLMOD   = 0x3A
	ST7789_RAMWRC   = 0x3C
	ST7789_RAMRDC   = 0x3E

	ST7789_TESCAN   = 0x44
	ST7789_RDTESCAN = 0x45

	ST7789_WRDISBV  = 0x51
	ST7789_RDDISBV  = 0x52
	ST7789_WRCTRLD  = 0x53
	ST7789_RDCTRLD  = 0x54
	ST7789_WRCACE   = 0x55
	ST7789_RDCABC   = 0x56
	ST7789_WRCABCMB = 0x5E
	ST7789_RDCABCMB = 0x5F

	ST7789_RDABCSDR = 0x68

	ST7789_RDID1 = 0xDA
	ST7789_RDID2 = 0xDB
	ST7789_RDID3 = 0xDC

	ST7789_RAMCTRL = 0xB0
	ST7789_RGBCTRL = 0xB1
	ST7789_PORCTRL = 0xB2
	ST7789_FRCTRL1 = 0xB3

	ST7789_GCTRL = 0xB7
	ST7789_DGMEN = 0xBA
	ST7789_VCOMS = 0xBB

	ST7789_LCMCTRL  = 0xC0
	ST7789_IDSET    = 0xC1
	ST7789_VDVVRHEN = 0xC2

	ST7789_VRHS     = 0xC3
	ST7789_VDVSET   = 0xC4
	ST7789_VCMOFSET = 0xC5
	ST7789_FRCTR2   = 0xC6
	ST7789_CABCCTRL = 0xC7
	ST7789_REGSEL1  = 0xC8
	ST7789_REGSEL2  = 0xCA
	ST7789_PWMFRSEL = 0xCC

	ST7789_PWCTRL1   = 0xD0
	ST7789_VAPVANEN  = 0xD2
	ST7789_CMD2EN    = 0xDF5A6902
	ST7789_PVGAMCTRL = 0xE0
	ST7789_NVGAMCTRL = 0xE1
	ST7789_DGMLUTR   = 0xE2
	ST7789_DGMLUTB   = 0xE3
	ST7789_GATECTRL  = 0xE4
	ST7789_PWCTRL2   = 0xE8
	ST7789_EQCTRL    = 0xE9
	ST7789_PROMCTRL  = 0xEC
	ST7789_PROMEN    = 0xFA
	ST7789_NVMSET    = 0xFC
	ST7789_PROMACT   = 0xFE

	// Colours for convenience
	ST7789_BLACK   = 0x0000 // 0b 00000 000000 00000
	ST7789_BLUE    = 0x001F // 0b 00000 000000 11111
	ST7789_GREEN   = 0x07E0 // 0b 00000 111111 00000
	ST7789_RED     = 0xF800 // 0b 11111 000000 00000
	ST7789_CYAN    = 0x07FF // 0b 00000 111111 11111
	ST7789_MAGENTA = 0xF81F // 0b 11111 000000 11111
	ST7789_YELLOW  = 0xFFE0 // 0b 11111 111111 00000
	ST7789_WHITE   = 0xFFFF // 0b 11111 111111 11111
)

type ST7789 struct {
	spi    SPI
	dc     PIN
	rst    PIN
	led    PIN
	width  int
	height int
	mux    sync.Mutex
}

func (s *ST7789) begin() {
	s.Reset()
	s.init()
}

// ExchangeData
//
//	@Description: 将数据写入SPI,isData为true表示写入的是数据,反之则是命令(非线程安全,请使用 Tx 包裹执行)
//	@receiver s
//	@param data 需要发送的数据
//	@param isData 是否是数据类型
func (s *ST7789) ExchangeData(isData bool, data []byte) {
	if isData {
		s.dc.High()
	} else {
		s.dc.Low()
	}
	s.spi.SpiTransmit(data)
}

// Command
//
//	@Description: 写入显示命令(非线程安全,请使用 Tx 包裹执行)
//	@receiver s
//	@param data 数据
func (s *ST7789) Command(data byte) {
	s.ExchangeData(false, []byte{data})
}

// Tx
//
//	@Description: 加锁执行，确保命令执行连续且原子
//	@receiver s
//	@param call 执行函数
func (s *ST7789) Tx(call func()) {
	s.mux.Lock()
	defer s.mux.Unlock()
	call()
}

// SendData
//
//	@Description: 写入显示数据(非线程安全,请使用 Tx 包裹执行)
//	@receiver s
//	@param data 数据
func (s *ST7789) SendData(data ...byte) {
	s.ExchangeData(true, data)
}

// Reset
//
//	@Description: Reset the display, if reset pin is connected.(线程安全)
//	@receiver s
func (s *ST7789) Reset() {
	s.Tx(func() {
		s.rst.High()
		time.Sleep(time.Millisecond * 100)
		s.rst.Low()
		time.Sleep(time.Millisecond * 100)
		s.rst.High()
		time.Sleep(time.Millisecond * 100)
	})
}

// init
//
//	@Description: Initialize the display. Broken out as a separate function so it can be overridden by other displays in the future.
//	@receiver s
func (s *ST7789) init() {
	s.Tx(func() {
		s.Command(0x11)
		time.Sleep(time.Millisecond * 150)

		s.Command(0x36)
		s.SendData(0x00)

		s.Command(0x3A)
		s.SendData(0x05)

		s.Command(0xB2)
		s.SendData(0x0C, 0x0C)

		s.Command(0xB7)
		s.SendData(0x35)

		s.Command(0xBB)
		s.SendData(0x1A)

		s.Command(0xC0)
		s.SendData(0x2C)

		s.Command(0xC2)
		s.SendData(0x01)

		s.Command(0xC3)
		s.SendData(0x0B)

		s.Command(0xC4)
		s.SendData(0x20)

		s.Command(0xC6)
		s.SendData(0x0F)

		s.Command(0xD0)
		s.SendData(0xA4, 0xA1)

		s.Command(0x21)

		s.Command(0xE0)
		s.SendData(
			0x00,
			0x19,
			0x1E,
			0x0A,
			0x09,
			0x15,
			0x3D,
			0x44,
			0x51,
			0x12,
			0x03,
			0x00,
			0x3F,
			0x3F)

		s.Command(0xE1)
		s.SendData(
			0x00,
			0x18,
			0x1E,
			0x0A,
			0x09,
			0x25,
			0x3F,
			0x43,
			0x52,
			0x33,
			0x03,
			0x00,
			0x3F,
			0x3F)
		s.Command(0x29)

		time.Sleep(time.Millisecond * 100) // 100 ms
	})
}

// setWindow
//
//	@Description: Set the pixel address window for proceeding drawing commands. x0 and
//	   x1 should define the minimum and maximum x pixel bounds.  y0 and y1
//	   should define the minimum and maximum y pixel bound.
//	@receiver s
//	@param x0 区域开始X轴位置(包含)
//	@param y0 区域开始Y轴位置(包含)
//	@param x1 区域结束X轴位置(包含)
//	@param y1 区域结束Y轴位置(包含)
func (s *ST7789) setWindow(x0, y0, x1, y1 int) {
	s.Command(ST7789_CASET) // Column addr set
	s.SendData(byte(
		x0>>8),
		byte(x0), // XSTART
		byte(x1>>8),
		byte(x1), // XEND
	)
	s.Command(ST7789_RASET) // Row addr set
	s.SendData(
		byte(y0>>8),
		byte(y0), // YSTART
		byte(y1>>8),
		byte(y1), // YEND
	)
	s.Command(ST7789_RAMWR) // write to RAM
}

// Flush
//
//	@Description: 将画布上的图像绘制到屏幕上(线程安全)
//	@receiver s
//	@param canvas 画布
func (s *ST7789) Flush(canvas *Canvas) {
	tmp := make([]byte, len(canvas.buffer))
	copy(tmp, canvas.buffer)
	s.Tx(func() {
		s.setWindow(canvas.x0, canvas.y0, canvas.x1, canvas.y1)
		s.ExchangeData(true, tmp)
	})
}

// Size
//
//	@Description: 获取显示器尺寸
//	@receiver s
//	@return *image.Point 尺寸
func (s *ST7789) Size() *image.Point {
	return &image.Point{
		X: s.width,
		Y: s.height,
	}
}

// GetFullScreenCanvas
//
//	@Description: 获取全屏画布
//	@receiver s
//	@return *Canvas 画布
func (s *ST7789) GetFullScreenCanvas() *Canvas {
	return &Canvas{
		device: s,
		x0:     0,
		y0:     0,
		x1:     s.width - 1,
		y1:     s.height - 1,
		width:  s.width,
		height: s.height,
		buffer: make([]byte, s.width*s.height*2),
	}
}

// GetCanvas
//
//	@Description: 获取画布
//	@receiver s
//	@param x0 区域X轴起始(包含)
//	@param y0 区域Y轴起始(包含)
//	@param x1 区域X轴截止(包含)
//	@param y1 区域X轴截止(包含)
//	@return *Canvas
func (s *ST7789) GetCanvas(x0, y0, x1, y1 int) *Canvas {
	width := x1 - x0 + 1
	height := y1 - y0 + 1
	return &Canvas{
		device: s,
		x0:     x0,
		y0:     y0,
		x1:     x1,
		y1:     y1,
		width:  width,
		height: height,
		buffer: make([]byte, width*height*2),
	}
}

// computeAlpha
//
//	@Description: 混合背景色
//	@param color 当前色值
//	@param bg 背景色值
//	@param alpha 当前色Alpha值
//	@param bgAlpha 背景色Alpha值
//	@return uint32
func computeAlpha(color uint32, bg uint16, alpha, bgAlpha uint32) uint32 {
	return (color*alpha + bgAlpha*uint32(bg)) / 255
}

// ColorToRgb565
//
//	@Description: 转换 color.Color 为RGB565
//	@param c 当前颜色
//	@param backgroundColor 背景颜色(RGB565)
//	@return uint16 RGB565色值
func ColorToRgb565(c color.Color, backgroundColor uint16) uint16 {
	r, g, b, a := c.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	a >>= 8
	f := 255 - a
	if f > 0 {
		br, bg, bb := Rgb565ToRgb(backgroundColor)
		r = computeAlpha(r, br, a, f)
		g = computeAlpha(g, bg, a, f)
		b = computeAlpha(b, bb, a, f)
	}
	return uint16(((r & 0xF8) << 8) | ((g & 0xFC) << 3) | (b >> 3))
}

// Rgb565ToRgb
//
//	@Description: 转换RGB565为标准RGB
//	@param c RGB565色值
//	@return r RGB(R色值)
//	@return g RGB(G色值)
//	@return b RGB(B色值)
func Rgb565ToRgb(c uint16) (r, g, b uint16) {
	r = (c >> 8) & 0xF8
	g = (c >> 3) & 0xFC
	b = (c & 0x1f) << 3
	return
}

type BaseCanvas interface {
	//
	// SetRGB565
	//  @Description: 设置指定坐标RGB565色值
	//  @param x X轴
	//  @param y Y轴
	//  @param c RBG565色值
	//
	SetRGB565(x, y int, c uint16)
	//
	// GetRGB565
	//  @Description: 获取指定坐标RGB565色值
	//  @param x X轴
	//  @param y Y轴
	//  @result RGB565色值
	//
	GetRGB565(x, y int) uint16
}

// Canvas
// @Description: 画布
type Canvas struct {
	device *ST7789
	x0     int    // X轴画布起始偏移
	y0     int    // Y轴画布起始偏移
	x1     int    // X轴画布结束偏移
	y1     int    // Y轴画布结束偏移
	width  int    // 画布宽度
	height int    // 画布高度
	buffer []byte // 缓冲区
}

// SetRGB565
//
//	@Description: 设置缓存区指定坐标的RGB565色值
//	@receiver d
//	@param x X轴坐标
//	@param y Y轴坐标
//	@param c RGB565色值
func (d *Canvas) SetRGB565(x, y int, c uint16) {
	index := d.getBufferBeginIndex(x, y)
	d.buffer[index] = byte(c >> 8)
	d.buffer[index+1] = byte(c)
}

// GetRGB565
//
//	@Description: 获取缓存区指定坐标的RGB565色值
//	@receiver d
//	@param x X轴坐标
//	@param y Y轴坐标
//	@return uint16 RGB565色值
func (d *Canvas) GetRGB565(x, y int) uint16 {
	index := d.getBufferBeginIndex(x, y)
	return (uint16(d.buffer[index]) << 8) + uint16(d.buffer[index+1])
}

// GetColor
//
//	@Description: 获取缓冲区指定坐标RGBA色值(由于该值从RBG565转换而来,故A值始终为1)
//	@receiver d
//	@param x X轴坐标
//	@param y Y轴坐标
//	@return color.Color
func (d *Canvas) GetColor(x, y int) color.Color {
	rgb565 := d.GetRGB565(x, y)
	r, g, b := Rgb565ToRgb(rgb565)
	return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
}

// SetColor
//
//	@Description: 设置缓冲区指定坐标的色值
//	@receiver d
//	@param x X轴坐标
//	@param y Y轴坐标
//	@param c 色值
func (d *Canvas) SetColor(x, y int, c color.Color) {
	value := ColorToRgb565(c, d.GetRGB565(x, y))
	d.SetRGB565(x, y, value)
}

// getBufferBeginIndex
//
//	@Description: 获取缓冲区
//	@receiver d
//	@param x X轴坐标
//	@param y Y轴坐标
//	@return int 缓冲区开始下标
func (d *Canvas) getBufferBeginIndex(x, y int) int {
	return (y*d.width + x) * 2
}

// Flush
//
//	@Description: 将缓冲区内容刷新到屏幕上
//	@receiver d
func (d *Canvas) Flush() {
	d.device.Flush(d)
}

// DrawImage
//
//	@Description: 将图像绘制到画布缓冲区中
//	@receiver d
//	@param img 图像
func (d *Canvas) DrawImage(img image.Image) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y && y < d.height; y++ {
		for x := bounds.Min.X; x < bounds.Max.X && x < d.width; x++ {
			d.SetColor(x, y, img.At(x, y))
		}
	}
}

// Clear
//
//	@Description: 清空画布缓冲区数据
//	@receiver d
func (d *Canvas) Clear() {
	for x := 0; x < d.width; x++ {
		for y := 0; y < d.height; y++ {
			d.SetRGB565(x, y, 0)
		}
	}
}

type SPI interface {
	//
	// SpiSpeed
	//  @Description: 设置SPI速率
	//  @param speed
	//
	SpiSpeed(speed uint32)
	//
	// SetSpiMode2
	//  @Description:设置为Mode 2 CPOL=1, CPHA=0模式
	//
	SetSpiMode2()
	//
	// SpiTransmit
	//  @Description: 发送数据
	//  @param data 需要发送的数据
	//
	SpiTransmit(data []byte)
}

type PIN interface {
	//
	// High
	//  @Description:输出为高电频
	//
	High()
	//
	// Low
	//  @Description:设置为低电频
	//
	Low()
	//
	// SetOutput
	//  @Description:设置为输出模式
	//
	SetOutput()
}

// NewST7789
//
//	@Description: ST7789显示驱动
//	@param spi SPI通信端口
//	@param dc 引脚DC
//	@param rst 引脚RES
//	@param led 引脚BLK
//	@param width 显示宽度
//	@param height 显示高度
//	@return *ST7789
//	@return error 创建失败
func NewST7789(spi SPI, dc, rst, led PIN, width, height int) *ST7789 {
	s := &ST7789{
		spi:    spi,
		dc:     dc,
		rst:    rst,
		led:    led,
		width:  width,
		height: height,
	}
	// Set DC as output.
	s.dc.SetOutput()
	// Setup reset as output
	s.rst.SetOutput()
	// Turn on the backlight LED
	s.led.SetOutput()
	s.led.High()
	// Set SPI to mode 0, MSB first.
	spi.SetSpiMode2()
	spi.SpiSpeed(SPI_CLOCK_HZ)
	s.begin()
	return s
}
