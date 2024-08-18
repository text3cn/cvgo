package routing

import (
	"client/internal/api"
	"cvgo/provider/httpserver"
)

// 路由定义
func Router(engine *httpserver.Engine) {
	engine.Get("/", api.Index)
}