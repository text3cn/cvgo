package clog

import (
	"fmt"
	"github.com/spf13/cast"
)

const (
	// 前景颜色
	black  = "\033[30m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	pink   = "\033[35m"
	cyan   = "\033[36m"
	gray   = "\033[37m"

	ColorBlack  = black
	ColorRed    = red
	ColorGreen  = green
	ColorYellow = yellow
	ColorBlue   = blue
	ColorPink   = pink
	ColorCyan   = cyan
	ColorGray   = gray

	// 背景颜色
	blackBG  = "\033[40m"
	redBG    = "\033[41m"
	greenBG  = "\033[42m"
	yellowBG = "\033[43m"
	blueBG   = "\033[44m"
	pinkBG   = "\033[45m"
	cyanBG   = "\033[46m"
	grayBG   = "\033[47m"

	// 文本样式
	textBlob           = "\033[1m" // 粗体
	textDownplay       = "\033[2m" // 淡化
	textItalic         = "\033[3m" // 斜体
	textUnderline      = "\033[4m" // 下划线
	textBlink          = "\033[5m" // 闪烁
	textReverseDisplay = "\033[7m" // 文本背景和前景颜色对调
	textInvisible      = "\033[7m" // 隐藏文本（文本不可见，但仍占用空间）

	// 重置所有颜色和样式: fmt.Printf(reset)
	reset = "\033[0m"

	// 组合效果
	// 这些转义码可以通过组合使用来实现不同的文本效果，
	// 例如，要创建红色粗体文本，可以使用 \033[31;1m，然后使用 \033[0m 重置文本样式。
	// fmt.Sprintf(pink + greenBG + "hello world" + reset)

	// 高亮加粗白
	whiteHighlight = "\033[1;97m"
)

// 红色
func RedPrintln(output ...any) {
	colorPirntln(red, output...)
}
func RedPrintf(str string, v ...any) {
	colorPrintf(red, str, v...)
}
func RedSprintf(str any, rest ...string) string {
	return colorSprintf(red, str, rest...)
}

// 紫色
func PinkPrintln(output ...any) {
	colorPirntln(pink, output...)
}
func PinkPrintf(str string, v ...any) {
	colorPrintf(pink, str, v...)
}
func PinkSprintf(str any, rest ...string) string {
	return colorSprintf(pink, str, rest...)
}

// 绿色
func GreenPrintln(output ...any) {
	colorPirntln(green, output...)
}
func GreenPrintf(str string, v ...any) {
	colorPrintf(green, str, v...)
}
func GreenSprintf(str any, rest ...string) string {
	return colorSprintf(green, str, rest...)
}

// 青色
func CyanPrintln(output ...any) {
	colorPirntln(cyan, output...)
}
func CyanPrintf(str string, v ...any) {
	colorPrintf(cyan, str, v...)
}
func CyanSprintf(str any, rest ...string) string {
	return colorSprintf(cyan, str, rest...)
}

// 蓝色
func BluePrintln(output ...any) {
	colorPirntln(blue, output...)
}
func BluePrintf(str string, v ...any) {
	colorPrintf(blue, str, v...)
}
func BlueSprintf(str any, rest ...string) string {
	return colorSprintf(blue, str, rest...)
}

// 黄色
func YellowPrintln(output ...any) {
	colorPirntln(yellow, output...)
}
func YellowPrintf(str string, v ...any) {
	colorPrintf(yellow, str, v...)
}
func YellowSprintf(str any, rest ...string) string {
	return colorSprintf(yellow, str, rest...)
}

func colorPirntln(color string, output ...any) {
	//fmt.Println(yellow+" %v "+reset, str)
	var str string
	if len(output) > 0 {
		for _, v := range output {
			str = str + cast.ToString(v) + " "
		}
	}
	fmt.Println(color + str + reset)
}
func colorPrintf(color string, str string, v ...any) {
	fmt.Printf(color+str+reset, v...)
}

func colorSprintf(color string, str any, rest ...string) string {
	restColor := reset
	if len(rest) > 0 {
		restColor = rest[0]
	}
	return fmt.Sprintf(color+" %v "+restColor, str)
}
