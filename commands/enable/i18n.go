package enable

import (
	"cvgo/kvs"
	"cvgo/kvs/kvsKey"
	"cvgo/paths"
	"github.com/textthree/cvgokit/filekit"
	"github.com/textthree/cvgokit/gokit"
	"github.com/textthree/provider/clog"
	"path/filepath"
)

// 添加 i18n 支持
func addI18n() {
	paths.CheckRunAtModuleRoot()

	modName, err := gokit.GetModuleName()
	if err != nil {
		panic(err)
	}
	kv := kvs.Instance()
	path := paths.NewModulePath()
	routingFile = path.RoutingGo()
	webFrameworkKey := kvsKey.ModuleWebFramework(modName)
	i18nStorageKey := kvsKey.ModuleI18n(modName)
	if val, _ := kv.GetBool(i18nStorageKey); val {
		log.Info("i18n 已经添加过了，无法重复执行。")
		return
	}

	frameworkType, _ := kv.GetString(webFrameworkKey)
	switch frameworkType {
	case "cvgo":
		cvgoAddI18n()
	case "fiber":
		fiberAddI18n(kv, path)
	}

	// 标识已添加 i18n
	kv.Set(i18nStorageKey, true)
	clog.GreenPrintln("添加 i18n 成功")
}

// cvgo
func cvgoAddI18n() {
	// instance.go 中声明变量
	workPath := paths.NewWorkPath()
	modulePath := paths.NewModulePath()
	instanceFile := workPath.InstanceGo()
	content := `    "github.com/textthree/provider/i18n"`
	filekit.AddContentUnderLine(instanceFile, "import (", content)

	content = "\n" + `var I18n i18n.Service`
	err := filekit.FileAppendContent(instanceFile, content)
	if err != nil {
		log.Error(err)
	}

	// boot/init.go 中获取实例
	initFile := modulePath.BootInitGo()
	content = "\n" + `
	"github.com/textthree/provider/i18n"`
	filekit.AddContentUnderLine(initFile, "import (", content)

	content = "    app.I18n = provider.Services.NewSingle(i18n.Name).(i18n.Service)"
	filekit.AddContentUnderLine(initFile, "func init() {", content)

	// 创建语言包 json 文件
	createJsonLanguagePackage()

	// 启用中间件
	content = `    cvgomiddleware "github.com/textthree/cvgoweb/middleware"`
	err = filekit.AddContentUnderLine(routingFile, "import (", content)
	if err != nil {
		log.Error("修改 routing.go 失败", err)
		return
	}
	content = `    engine.UseMiddleware(cvgomiddleware.I18n())`
	filekit.AddContentUnderLine(routingFile, "func Routes(engine *httpserver.Engine) {", content)
}

func fiberAddI18n(kv kvs.KvStorage, modPath *paths.ModulePath) {
	// instance.go 中声明变量
	workPath := paths.NewWorkPath()
	modulePath := paths.NewModulePath()
	instanceFile := workPath.InstanceGo()
	content := `    "github.com/textthree/provider/i18n"`
	filekit.AddContentUnderLine(instanceFile, "import (", content)

	content = "\n" + `var I18n i18n.Service`
	err := filekit.FileAppendContent(instanceFile, content)
	if err != nil {
		log.Error(err)
	}

	// boot -> init.go 中获取实例
	initFile := modulePath.BootInitGo()
	content = "\n" + `
	"github.com/textthree/provider/i18n"`
	filekit.AddContentUnderLine(initFile, "import (", content)

	content = "    app.I18n = provider.Services.NewSingle(i18n.Name).(i18n.Service)"
	filekit.AddContentUnderLine(initFile, "func init() {", content)

	// 创建语言包 json 文件
	createJsonLanguagePackage()

	// Create i18n middleware
	content = `package middleware

import (
	"` + kv.GetWorkspaceName() + `/app"
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
`
	err = filekit.FilePutContents(modPath.I18nMiddlewareGo(), content)
	if err != nil {
		panic(err)
	}

	// 启用中间件
	content = `    app.Use(middleware.I18n)`
	filekit.AddContentUnderLine(routingFile, "func Routes(app *fiber.App) {", content)
}

func createJsonLanguagePackage() {
	content := `{
  "hello": "你好",
  "world": {
    "china": "中国"
  }
}`
	filekit.FilePutContents(filepath.Join(pwd, "i18n", "zh.json"), content)
	content = `{
  "hello": "hello",
  "world": {
    "china": "china"
  }
}`
	filekit.FilePutContents(filepath.Join(pwd, "i18n", "en.json"), content)
}
