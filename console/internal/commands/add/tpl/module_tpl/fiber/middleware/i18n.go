package middleware

import (
	"cvgo/app"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"os"
)

func I18n(ctx *fiber.Ctx) error {
	lngCode := ctx.Get("Language", "en")
	app.I18n.SetLngCode(lngCode)
	if !app.I18n.LoadedPackage(lngCode) {
		lngpkg := viper.New()
		pwd, _ := os.Getwd()
		app.Log.Trace("加载语言包:", pwd+string(os.PathSeparator)+"i18n"+string(os.PathSeparator)+lngCode+".json")
		lngpkg.AddConfigPath(pwd + string(os.PathSeparator) + "i18n")
		lngpkg.SetConfigName(lngCode)
		lngpkg.SetConfigType("json")
		if err := lngpkg.ReadInConfig(); err != nil {
			app.Log.Error(err, lngCode)
		}
		app.I18n.SetLanguagePackage(lngCode, lngpkg)
	}
	return ctx.Next()
}
