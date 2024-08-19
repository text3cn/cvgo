package middleware

import (
	"cvgo/app/common/dto"
	"cvgo/errcode"
	"cvgo/kit/cryptokit"
	"cvgo/kit/jsonkit"
	"cvgo/provider"
	"cvgo/provider/clog"
	"cvgo/provider/config"
	"cvgo/provider/httpserver"
	"errors"
	"github.com/spf13/cast"
)

func Auth() httpserver.MiddlewareHandler {
	cfg := provider.Services.NewSingle(config.Name).(config.Service)
	secret := cfg.GetTokenSecret()

	return func(context *httpserver.Context) error {
		token, _ := context.Req.Header("Authorization")
		var uid string
		if token != "" {
			uid = cryptokit.DynamicDecrypt(secret, token)
		}
		clog.PinkPrintf("tokenKey=%s, token=%s, userId=%s \n", secret, token, uid)
		if token == "" || cast.ToInt64(uid) == 0 {
			ret := dto.BaseRes{
				ApiCode:    errcode.AuthorizationFailed,
				ApiMessage: "Authorization failed.",
			}
			info := string(jsonkit.JsonEncode(ret))
			return errors.New(info)
		}
		context.SetVal("uid", uid)
		context.Next()
		return nil
	}
}
