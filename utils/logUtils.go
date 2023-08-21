package utils

import (
	"fmt"
	"log"
	"os"
)

const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
)

var myLogger *log.Logger

func init() {
	myLogger = log.New(os.Stdout, "[Default]", log.Lshortfile|log.Ldate|log.Ltime)
}

// 重写 log 的Println 方法，修改调用堆栈的追踪深度，以便调试
func overridePrintln(l *log.Logger, v ...any) {
	err := l.Output(3, fmt.Sprintln(v...))
	if err != nil {
		return
	}
}

func Info(v ...any) {
	myLogger.SetPrefix(colorGreen + "[Info]" + colorReset)
	overridePrintln(myLogger, colorGreen+fmt.Sprint(v...)+colorReset)
}
func Error(v ...any) {
	myLogger.SetPrefix(colorRed + "[Error]" + colorReset)
	overridePrintln(myLogger, colorRed+fmt.Sprint(v...)+colorReset)
	os.Exit(1)
}
func Warning(v ...any) {
	myLogger.SetPrefix(colorYellow + "[Warning]" + colorReset)
	overridePrintln(myLogger, colorYellow+fmt.Sprint(v...)+colorReset)
	// 上一行参数传v表示整体当成数组传入，参数传v... 表示多个参数分别传入
}

// InfoF 带格式化的信息日志
func InfoF(format string, v ...any) {
	myLogger.SetPrefix(colorGreen + "[Info]" + colorReset)
	overridePrintln(myLogger, colorGreen+fmt.Sprintf(format, v...)+colorReset)
}

// ErrorF 带格式化的错误日志
func ErrorF(format string, v ...any) {
	myLogger.SetPrefix(colorRed + "[Error]" + colorReset)
	overridePrintln(myLogger, colorRed+fmt.Sprintf(format, v...)+colorReset)
	os.Exit(1)
}

// WarningF 带格式化的警告日志
func WarningF(format string, v ...any) {
	myLogger.SetPrefix(colorYellow + "[Warning]" + colorReset)
	overridePrintln(myLogger, colorYellow+fmt.Sprintf(format, v...)+colorReset)
}
