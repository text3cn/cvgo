package clog

import (
	"fmt"
	"github.com/spf13/cast"
	"log"
	"os"
	"strconv"
	"strings"
)

func (self *ClogService) output(level byte, v ...interface{}) error {
	// calldepth 层次说明：本函数（第一层） -> service.go（第二层） -> 用户调用时（第三层）
	// 最后一个参数可以传 deep:number 来设置层数，例如：deep:5
	calldepth := 3
	length := len(v)
	last := cast.ToString(v[length-1])
	if strings.HasPrefix(last, "deep:") {
		arr := strings.Split(last, ":")
		calldepth, _ = strconv.Atoi(arr[1])
		length--
	}
	// 格式化输出样式
	str := ""
	formatStr := "%v"
	for i := 0; i < length; i++ {
		str += fmt.Sprintf(formatStr, v[i]) + " "
	}
	switch level {
	case trace:
		formatStr = "\033[37m[TRACE] %s\033[0m"
	case debug:
		formatStr = "\033[32m[DEBUG] %s\033[0m"
	case info:
		formatStr = "\033[36m[INFO] %s\033[0m"
	case warn:
		formatStr = "\033[33m[WARN] %s\033[0m"
	case err:
		formatStr = "\033[31m[ERROR] %s\033[0m"
	case fatal:
		formatStr = "\033[31m[FATAL] %s\033[0m"
	}
	str = fmt.Sprintf(formatStr, str)
	//	s := fmt.Sprintf(formatStr, v...)
	// fmt.Printf("%c[%d;%d;%dm%s%c[0m", 0x1B, 1, 97, 31, message, 0x1B)
	//s = fmt.Sprintf("%c[%d;%d;%dm%s%c[0m", 0x1B, 1, 97, 31, v, 0x1B)

	_log := log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)
	return _log.Output(calldepth, str)
}
