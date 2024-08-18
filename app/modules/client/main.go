package main

import (
	_ "client/internal/boot"
	"client/internal/routing"
	"cvgo/provider/httpserver/cvgohttp"
)

// @title client API Doc
// @description 使用 [swaggo](https://github.com/swaggo/swag/blob/master/README_zh-CN.md) 编写。
// @host localhost:9000
// @BasePath /
func main() {

	cvgohttp.Run(routing.Router)
}
