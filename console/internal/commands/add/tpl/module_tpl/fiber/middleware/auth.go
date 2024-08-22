package middleware

import (
	"cvgo/app/common/dto"
	"cvgo/app"
	"cvgo/kit/cryptokit"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
)

// 获取请求头处理
func AuthMiddleware(ctx *fiber.Ctx) error {
	app.Log.Trace(ctx.Request())
	auth := ctx.Get("Authorization")
	if auth == "" {
		return ctx.JSON(dto.BaseRes{
			ApiCode:    1,
			ApiMessage: "No Login",
		})
	}
	uid := cryptokit.DynamicDecrypt("Token 密钥", auth)
	ctx.Locals("uid", uid)
	if cast.ToInt64(uid) == 0 {
		return ctx.JSON(dto.BaseRes{
			ApiCode:    1,
			ApiMessage: "Login invalid",
		})
	}
	// 继续处理请求
	return ctx.Next()
}
