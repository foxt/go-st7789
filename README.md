# ST7789

使用Goland实现的操作ST7789,适用于无CS引脚的240x204 LCD显示屏。 目前仅在Raspberry zero 2w上测试通过。

本库根据Python版 https://github.com/solinnovay/Python_ST7789 移植而来，并在此基础上实现了RGBA转RGB565,支持透明图层。
# 安装
```shell
go get github.com/manx98/go-st7789
```
使用示例位于 sample/sample.go
# 感谢
1. Python原始实现 https://github.com/solinnovay/Python_ST7789
2. GPIO库 https://github.com/stianeikeland/go-rpio/
