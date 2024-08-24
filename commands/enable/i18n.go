package enable

import (
	"cvgo/console"
	"cvgo/console/internal/paths"
	"github.com/textthree/cvgokit/filekit"
	"github.com/textthree/cvgokit/gokit"
	"github.com/textthree/provider/clog"
	"path/filepath"
)

// 添加 i18n 支持
func addI18n() {
	if !paths.CheckRunAtModuleRoot() {
		return
	}
	modName, err := gokit.GetModuleName()
	if err != nil {
		panic(err)
	}
	routingFile = filepath.Join(pwd, "internal", "routing", "routing.go")
	kv := console.NewKvStorage(filekit.GetParentDir(3))
	webFrameworkKey := modName + "." + "webframework"
	i18nStorageKey := modName + "." + "i18n"
	if val, _ := kv.GetBool(i18nStorageKey); val {
		log.Info("i18n 已经添加过了，无法重复执行。")
		return
	}

	// instance.go 中声明变量
	instanceFile := filepath.Join(filekit.GetParentDir(2), "instance.go")
	content := `    "github.com/textthree/provider/i18n"`
	filekit.AddContentUnderLine(instanceFile, "import (", content)

	content = "\n" + `var I18n i18n.Service`
	err = filekit.FileAppendContent(instanceFile, content)
	if err != nil {
		log.Error(err)
	}

	// boot/init.go 中获取实例
	initFile := filepath.Join(pwd, "internal", "boot", "init.go")
	content = "\n" + `
	"github.com/textthree/provider/i18n"
`
	filekit.AddContentUnderLine(initFile, "import (", content)

	content = "    app.I18n = provider.Services.NewSingle(i18n.Name).(i18n.Service)"
	filekit.AddContentUnderLine(initFile, "func init() {", content)

	// 创建语言包 json 文件
	content = `{
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

	// 启用中间件
	frameworkType, _ := kv.GetString(webFrameworkKey)
	switch frameworkType {
	case "cvgo":
		useCvgoI18nMiddleware()
	case "fiber":
		useFiberI18nMiddleware()
	}

	// 标识已添加 i18n
	kv.Set(i18nStorageKey, true)
	clog.GreenPrintln("添加 i18n 成功")
}

// cvgo 框架启用 i18n 中间件
func useCvgoI18nMiddleware() {
	content := `    "github.com/textthree/provider/httpserver/middleware"
`
	err := filekit.AddContentUnderLine(routingFile, "import (", content)
	if err != nil {
		log.Error("修改 routing.go 失败", err)
		return
	}
	content = `
    engine.UseMiddleware(middleware.I18n())
`
	filekit.AddContentUnderLine(routingFile, "func Routes(engine *httpserver.Engine) {", content)
}

// fiber 框架添加 i18n 中间件
func useFiberI18nMiddleware() {
	// 拷贝自定义中间件模板
	src := filepath.Join(paths.FiberTplForModule(), "middleware", "i18n.go")
	dest := filepath.Join(pwd, "internal", "middleware", "i18n.go")
	filekit.CopyFile(src, dest)
	// 在路由中启用中间件
	content := `    app.Use(middleware.I18n)`
	filekit.AddContentUnderLine(routingFile, "func Routes(app *fiber.App) {", content)
}
