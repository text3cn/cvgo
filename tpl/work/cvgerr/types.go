// 扩展 errors 包，自定义错误类型
package cvgerr

import (
	"fmt"
)

var AllErrors = make(map[int]string)

type ApiError struct {
	Code    int
	Message string
}

func NewApiError(code int, message string) ApiError {
	if _, ok := AllErrors[code]; ok {
		msg := fmt.Sprintf("Duplicated error code=%d", code)
		panic(msg)
	}
	AllErrors[code] = message
	return ApiError{code, message}
}

// 请求参数解析错误
func ParseRequestParamsError() ApiError {
	return ApiError{
		Code:    ParseRequestParamsFailed.Code,
		Message: ParseRequestParamsFailed.Message,
	}
}
