package routing

import (
	"chord/internal/api"
	"chord/internal/middleware"
	"github.com/textthree/cvgoweb"
)

// 路由定义
func Routes(engine *httpserver.Engine) {
	engine.Cross()
	engine.Get("/", api.Index)

	// 需要验证 token 的路由分组
	authGroup := engine.Prefix("/api").UseMiddleware(middleware.Auth())
	{
		// cvgflag=authGroup
        authGroup.Get("user/info", api.GetUserInfo)
		authGroup.Get("/demo", api.Index)
	}
}
