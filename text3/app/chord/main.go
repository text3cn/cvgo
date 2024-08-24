package main

import (
	_ "chord/internal/boot"
	"chord/internal/routing"
	"github.com/textthree/cvgoweb/cvgohttp"
)


// @title chord API Doc
// @description 使用 [swaggo](https://github.com/swaggo/swag/blob/master/README_zh-CN.md) 编写。
// @host localhost:9003
// @BasePath /
func main() {
	cvgohttp.Run(routing.Routes)
}
