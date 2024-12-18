package ST7789

import (
	"fmt"
	"image"
	"image/color"
	"time"
)

// ColorMode 色彩模式
type ColorMode uint8

// ScreenType  显示屏类型
type ScreenType uint8

const (
	SPI_CLOCK_HZ = 40000000 // 40 MHz

	// Constants for interacting Width display registers.

	ST7789_NOP       = 0x00
	ST7789_SWRESET   = 0x01 // Software reset
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

	ST7789_POWSAVE    = 0xbc
	ST7789_DLPOFFSAVE = 0xbd

	// Colours for convenience
	ST7789_BLACK   = 0x0000 // 0b 00000 000000 00000
	ST7789_BLUE    = 0x001F // 0b 00000 000000 11111
	ST7789_GREEN   = 0x07E0 // 0b 00000 111111 00000
	ST7789_RED     = 0xF800 // 0b 11111 000000 00000
	ST7789_CYAN    = 0x07FF // 0b 00000 111111 11111
	ST7789_MAGENTA = 0xF81F // 0b 11111 000000 11111
	ST7789_YELLOW  = 0xFFE0 // 0b 11111 111111 00000
	ST7789_WHITE   = 0xFFFF // 0b 11111 111111 11111

	COLOR_MODE_65K   = ColorMode(0x50)
	COLOR_MODE_262K  = ColorMode(0x60)
	COLOR_MODE_12BIT = ColorMode(0x03)
	COLOR_MODE_16BIT = ColorMode(0x05)
	COLOR_MODE_18BIT = ColorMode(0x06)
	COLOR_MODE_16M   = ColorMode(0x07)

	// Screen320X240 Width 320,Height 240
	Screen320X240 = ScreenType(0)
	// Screen240X240 Width 240,Height 240
	Screen240X240 = ScreenType(1)
	// Screen135X240 Width 135,Height 240
	Screen135X240 = ScreenType(2)
)

// MADCTL ROTATIONS[rotation % 4]
var rotations = []byte{0x00, 0x60, 0xc0, 0xa0}

var width320 = [][]int{
	{240, 320, 0, 0},
	{320, 240, 0, 0},
	{240, 320, 0, 0},
	{320, 240, 0, 0},
}

var width240 = [][]int{
	{240, 240, 0, 0},
	{240, 240, 0, 0},
	{240, 240, 0, 80},
	{240, 240, 80, 0},
}

var width135 = [][]int{
	{135, 240, 52, 40},
	{240, 135, 40, 53},
	{135, 240, 53, 40},
	{240, 135, 40, 52},
}

type ST7789 struct {
	spi         SPI
	dc          PIN
	rst         PIN
	led         PIN
	width       int
	height      int
	xStart      int
	yStart      int
	rotationMap [][]int
}

// begin
//
//	@Description: 初始化
//	@receiver s
func (s *ST7789) begin() {
	s.HardReset()
	s.init()
}

// ExchangeData
//
//	@Description: Write data to SPI. If isData is true, it means that the data is written, otherwise it is a command (not thread-safe, please use the Tx package to execute)
//	@receiver s
//	@param data Data to be sent
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
//	@Description: Write display command (not thread-safe, please use Tx package to execute)
//	@receiver s
//	@param data 数据
func (s *ST7789) Command(data byte) {
	s.ExchangeData(false, []byte{data})
}

// SendData
//
//	@Description: Write display data
//	@receiver s
//	@param data 数据
func (s *ST7789) SendData(data ...byte) {
	s.ExchangeData(true, data)
}

// HardReset
//
//	@Description: 硬重启设备
//	@receiver s
func (s *ST7789) HardReset() {
	s.rst.High()
	time.Sleep(time.Millisecond * 100)
	s.rst.Low()
	time.Sleep(time.Millisecond * 100)
	s.rst.High()
	time.Sleep(time.Millisecond * 100)
}

// SoftReset
//
//	@Description: 软复位
//	@receiver s
func (s *ST7789) SoftReset() {
	s.Command(ST7789_SWRESET)
}

func (s *ST7789) init() {
	s.SleepMode(false)
	time.Sleep(time.Millisecond * 150)

	s.Rotation(0)

	s.ColorMode(COLOR_MODE_65K | COLOR_MODE_16BIT)

	s.Command(ST7789_PORCTRL)
	s.SendData(0x0C, 0x0C)

	s.Command(ST7789_GCTRL)
	s.SendData(0x35)

	s.Command(ST7789_VCOMS)
	s.SendData(0x1A)

	s.Command(ST7789_LCMCTRL)
	s.SendData(0x2C)

	s.Command(ST7789_VDVVRHEN)
	s.SendData(0x01)

	s.Command(ST7789_VRHS)
	s.SendData(0x0B)

	s.Command(ST7789_VDVSET)
	s.SendData(0x20)

	s.Command(ST7789_FRCTR2)
	s.SendData(0x0F)

	s.Command(ST7789_PWCTRL1)
	s.SendData(0xA4, 0xA1)

	s.InversionMode(true)

	s.Command(ST7789_PVGAMCTRL)
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

	s.Command(ST7789_NVGAMCTRL)
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
	s.Command(ST7789_DISPON)

	time.Sleep(time.Millisecond * 100) // 100 ms
}

// SetWindow
//
//	@Description: Set the pixel address window for proceeding drawing commands. X0 and
//	   X1 should define the minimum and maximum x pixel bounds.  Y0 and Y1
//	   should define the minimum and maximum y pixel bound.
//	@receiver s
//	@param X0 Region start X-axis position (inclusive)
//	@param Y0 Region start Y-axis position (inclusive)
//	@param X1 Region end X-axis position (inclusive)
//	@param Y1 Region end Y-axis position (inclusive)
func (s *ST7789) SetWindow(x0, y0, x1, y1 int) {
	s.Command(ST7789_CASET) // Column addr set
	x0 += s.xStart
	x1 += s.xStart
	s.SendData(byte(
		x0>>8),
		byte(x0), // XSTART
		byte(x1>>8),
		byte(x1), // XEND
	)
	s.Command(ST7789_RASET) // Row addr set
	y0 += s.yStart
	y1 += s.yStart
	s.SendData(
		byte(y0>>8),
		byte(y0), // YSTART
		byte(y1>>8),
		byte(y1), // YEND
	)
	s.Command(ST7789_RAMWR) // write to RAM
}

// FlushBitBuffer
//
//	@Description: Draw an image from the canvas to the screen
//	@receiver s
//	@param X0 Start position of the area on the X-axis (inclusive)
//	@param Y0 Start position of the area on the Y-axis (inclusive)
//	@param X1 End position of the area on the X-axis (inclusive)
//	@param Y1 End position of the area on the Y-axis (inclusive)
//	@param Buffer RGB565 image data
func (s *ST7789) FlushBitBuffer(x0, y0, x1, y1 int, buffer []byte) {
	s.SetWindow(x0, y0, x1, y1)
	s.ExchangeData(true, buffer)
}

// Size
//
//	@Description: Get the size of the display
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
		X0:     0,
		Y0:     0,
		X1:     s.width - 1,
		Y1:     s.height - 1,
		Width:  s.width,
		Height: s.height,
		Buffer: make([]byte, s.width*s.height*2),
	}
}

// GetCanvas
//
//	@Description: 获取画布
//	@receiver s
//	@param X0 Region start X-axis position (inclusive)
//	@param Y0 Region start Y-axis position (inclusive)
//	@param X1 Region end X-axis position (inclusive)
//	@param Y1 Region end Y-axis position (inclusive)
//	@return *Canvas
func (s *ST7789) GetCanvas(x0, y0, x1, y1 int) *Canvas {
	width := x1 - x0 + 1
	height := y1 - y0 + 1
	return &Canvas{
		device: s,
		X0:     x0,
		Y0:     y0,
		X1:     x1,
		Y1:     y1,
		Width:  width,
		Height: height,
		Buffer: make([]byte, width*height*2),
	}
}

// SleepMode
//
//	@Description: 是否启用显示休眠模式
//	@receiver s
//	@param enable 是否开启
func (s *ST7789) SleepMode(enable bool) {
	if enable {
		s.Command(ST7789_SLPIN)
	} else {
		s.Command(ST7789_SLPOUT)
	}
}

// InversionMode
//
//	@Description: 是否启用显示反转模式
//	@receiver s
//	@param enable 是否启用
func (s *ST7789) InversionMode(enable bool) {
	if enable {
		s.Command(ST7789_INVON)
	} else {
		s.Command(ST7789_INVOFF)
	}
}

// ColorMode
//
//	@Description: 设置颜色模式
//	@receiver s
//	@param mode
//		COLOR_MODE_65K, COLOR_MODE_262K, COLOR_MODE_12BIT,
//		COLOR_MODE_16BIT, COLOR_MODE_18BIT, COLOR_MODE_16M
func (s *ST7789) ColorMode(mode ColorMode) {
	s.Command(ST7789_COLMOD)
	s.SendData(uint8(mode) & 0x77)
}

// Rotation
//
//	@Description: 设置显示旋转
//	@receiver s
//	@param rotation
//	  	0-Portrait
//	  	1-Landscape
//	  	2-Inverted Portrait
//	  	3-Inverted Landscape
func (s *ST7789) Rotation(rotation uint8) {
	rotation = rotation % 4
	s.width = s.rotationMap[rotation][0]
	s.height = s.rotationMap[rotation][1]
	s.xStart = s.rotationMap[rotation][2]
	s.yStart = s.rotationMap[rotation][3]
	s.Command(ST7789_MADCTL)
	s.SendData(rotations[rotation%4])
}

// PowerSave
//
//	@Description:
//	@receiver s
//	@param mode
//		0 - off
//		1 - idle
//		2 - normal
//		4 - display off
func (s *ST7789) PowerSave(mode uint8) {
	if mode == 0 {
		s.Command(ST7789_POWSAVE)
		s.SendData(0xec | 3)
		s.Command(ST7789_DLPOFFSAVE)
		s.SendData(0xff)
		return
	}
	var is byte
	if mode&1 == 0 {
		is = 1
	}
	var ns byte
	if mode&2 == 0 {
		ns = 2
	}
	s.Command(ST7789_POWSAVE)
	s.SendData(0xec | ns | is)
	if mode&4 > 0 {
		s.Command(ST7789_DLPOFFSAVE)
		s.SendData(0xfe)
	}
}

// computeAlpha
//
//	@Description: Mix two colors together
//	@param color Foreground color value
//	@param bg Background color value
//	@param alpha Foreground alpha value
//	@param bgAlpha Background alpha value
//	@return uint32
func computeAlpha(color uint32, bg uint16, alpha, bgAlpha uint32) uint32 {
	return (color*alpha + bgAlpha*uint32(bg)) / 255
}

// ColorToRgb565
//
//	@Description: Convert color.Color to RGB565
//	@param c 当前颜色
//	@param backgroundColor background color (RGB565)
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
	//  @Description: Set the RGB565 color value of the specified coordinates
	//  @param x X轴
	//  @param y Y轴
	//  @param c RBG565色值
	//
	SetRGB565(x, y int, c uint16)
	//
	// GetRGB565
	//  @Description: Get the RGB565 color value of the specified coordinates
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
	X0 int //X-axis canvas starting offset
	Y0 int // Y-axis canvas starting offset
	X1 int //X-axis canvas end offset
	Y1 int // Y-axis canvas end offset
	Width int // Canvas width
	Height int // Canvas height
	Buffer []byte // buffer
}

// SetRGB565
//
//	@Description: Set the RGB565 color value of the specified coordinates in the buffer area
//	@receiver d
//	@param x X轴坐标
//	@param y Y轴坐标
//	@param c RGB565色值
func (d *Canvas) SetRGB565(x, y int, c uint16) {
	index := d.getBufferBeginIndex(x, y)
	d.Buffer[index] = byte(c >> 8)
	d.Buffer[index+1] = byte(c)
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
	return (uint16(d.Buffer[index]) << 8) + uint16(d.Buffer[index+1])
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
//	@Description: Returns the index of a pixel in the buffer
//	@receiver d
//	@param x X轴坐标
//	@param y Y轴坐标
//	@return int 缓冲区开始下标
func (d *Canvas) getBufferBeginIndex(x, y int) int {
	return (y*d.Width + x) * 2
}

// Flush
//
//	@Description: 将缓冲区内容刷新到屏幕上
//	@receiver d
func (d *Canvas) Flush() {
	d.FlushDirectly(d.Buffer)
}

// DrawImage
//
//	@Description: 将图像绘制到画布缓冲区中
//	@receiver d
//	@param img 图像
func (d *Canvas) DrawImage(img image.Image) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y && y < d.Height; y++ {
		for x := bounds.Min.X; x < bounds.Max.X && x < d.Width; x++ {
			d.SetColor(x, y, img.At(x, y))
		}
	}
}

// FlushDirectly
//
//	@Description: Draw the buffer content directly to the display area corresponding to the canvas. This method will not overwrite the canvas buffer.
//	@receiver d
//	@param buffer
func (d *Canvas) FlushDirectly(buffer []byte) {
	d.device.FlushBitBuffer(d.X0, d.Y0, d.X1, d.Y1, buffer)
}

// Clear
//
//	@Description: 清空画布缓冲区数据
//	@receiver d
func (d *Canvas) Clear() {
	for x := 0; x < d.Width; x++ {
		for y := 0; y < d.Height; y++ {
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
	// SetSpiMode3
	//  @Description: Set to Mode3 CPOL=1, CPHA=1 mode
	//
	SetSpiMode3()
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
// @Description: ST7789 display driver
// @param spi SPI communication port
// @param dc pin DC
// @param rst pin RES
// @param led pin BLK
// @param screen display type
// @return *ST7789
// @return error Creation failed
func NewST7789(spi SPI, dc, rst, led PIN, screen ScreenType) *ST7789 {
	s := &ST7789{
		spi: spi,
		dc:  dc,
		rst: rst,
		led: led,
	}
	switch screen {
	case Screen135X240:
		s.width = 135
		s.height = 240
		s.rotationMap = width135
	case Screen240X240:
		s.width = 240
		s.height = 240
		s.rotationMap = width240
	case Screen320X240:
		s.width = 320
		s.height = 240
		s.rotationMap = width320
	default:
		panic(fmt.Sprintf("Unsupported display. 320x240, 240x240 and 135x240 are supported."))
	}
	// Set DC as output.
	s.dc.SetOutput()
	// Setup reset as output
	s.rst.SetOutput()
	// Turn on the backlight LED
	s.led.SetOutput()
	s.led.High()
	spi.SpiSpeed(SPI_CLOCK_HZ)
	spi.SetSpiMode3()
	s.begin()
	return s
}
