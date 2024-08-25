package genfiber

import (
	"cvgo/config"
	"cvgo/kvs"
	"cvgo/kvs/kvsKey"
	"cvgo/paths"
	"github.com/textthree/cvgokit/filekit"
	"github.com/textthree/provider/clog"
	"path/filepath"
)

// 使用 fiber http 服务
func CreateWebserver(modName string, supportSwagger bool) {
	path := paths.NewModulePath(modName)
	moduleDir := path.ModulePath
	workspaceName, err := kvs.Instance().GetString(kvsKey.WorkspaceName)

	// 创建 index.api.go
	content := `package api

import "github.com/gofiber/fiber/v2"

func Index(ctx *fiber.Ctx) error {
    return ctx.SendString("Index")
}
`
	filekit.FilePutContents(path.IndeApiGo(), content)
	err = filekit.DeleteFile(path.IndeApiGitkeep())
	if err != nil {
		clog.RedPrintln(err.Error())
	}

	// boot init
	genBootInit(workspaceName, path.BootInitGo(), path.BootInitGitkeep())

	// 创建路由
	content = `package routing

import (
	"` + modName + `/internal/api"
	"` + modName + `/internal/middleware"
    "github.com/gofiber/fiber/v2"
)

// 路由定义
func Routes(app *fiber.App) {
	app.Get("/", api.Index)

	// 需要鉴权访问的路由组
	authGroup := app.Group("/api").Use(middleware.AuthMiddleware)
	{
		authGroup.Get("/demo", api.Index)
	}
}`
	filekit.FilePutContents(path.RoutingGo(), content)

	// 鉴权中间件
	genMiddleWare(path, workspaceName)

	// 创建 main.go
	swagger := ""
	if supportSwagger {
		genBootSwagger(path.BootSwaggerGo(), workspaceName)
		swagger = `
    // swagger 文档
	boot.SwaggerDoc(app)
`
	}
	mainGoFilePath := filepath.Join(moduleDir, "main.go")
	content = `package main

import (
	"` + modName + `/internal/boot"
	"` + modName + `/internal/routing"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)

func main() {
	app := fiber.New()

	// 跨域支持
	app.Use(cors.New())

	// panic 恢复
	boot.Recover(app)
	` + swagger + `
	// 路由
	routing.Routes(app)

	// 启动 http 服务
	log.Fatal(app.Listen(":" + boot.HttpServerPort))
}`
	filekit.FilePutContents(mainGoFilePath, content)

	// go.mod
	createGoModFile(moduleDir, modName)
}

func genBootSwagger(filePath, workspaceName string) {
	content := `package boot

import (
	"` + workspaceName + `/app"

	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"os"
	"path/filepath"
)

// Swagger
func SwaggerDoc(fiberApp *fiber.App) {
	env := app.Env()
	port := HttpServerPort
	var domain string
	if env == "" {
		domain = "127.0.0.1"
	} else if env == "alpha" {
		domain = "xxxx"
		port = "80"
	}
	// swagger-doc
	docFilePath := app.Config.GetSwagger().FilePath
	docFileName := filepath.Base(docFilePath)

	fiberApp.Get("/"+docFileName, func(c *fiber.Ctx) error {
		data, _ := os.ReadFile(docFilePath)
		return c.SendString(string(data))
	})
	url := "http://" + domain + ":" + port + "/" + docFileName
	fiberApp.Get("/swagger/*", swagger.New(swagger.Config{URL: url}))
	fmt.Println("  Swagger Doc: http://" + domain + ":" + HttpServerPort + "/swagger")
}
`
	filekit.FilePutContents(filePath, content)
}

func genMiddleWare(path *paths.ModulePath, workspaceName string) {
	// Create auth middleware
	content := `package middleware

import (
	"` + workspaceName + `/app"
	"` + workspaceName + `/common/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"github.com/textthree/cvgokit/cryptokit"
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
`
	err := filekit.FilePutContents(path.AuthMiddlewareGo(), content)
	if err != nil {
		panic(err)
	}

	// delete .gitkeep
	filekit.DeleteFile(path.MiddlewareGitkeep())
}

// 创建 go.mod 文件
func createGoModFile(modDir, modName string) {
	filePath := filepath.Join(modDir, "go.mod")
	content := `module ` + modName + `

go ` + config.GoVersion

	content += `
require (
	github.com/gofiber/fiber/v2 v2.52.5
	github.com/gofiber/swagger v1.1.0
	github.com/spf13/viper v1.19.0
	github.com/swaggo/swag v1.16.3
	github.com/textthree/provider v1.0.2
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/VictoriaMetrics/fastcache v1.12.2 // indirect
	github.com/allegro/bigcache v1.2.1 // indirect
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/coocood/freecache v1.2.4 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/spec v0.20.4 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.17.2 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/redis/go-redis/v9 v9.6.1 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.7.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/swaggo/files/v2 v2.0.0 // indirect
	github.com/textthree/cvgokit v1.0.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.51.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
	golang.org/x/net v0.23.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.13.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/mysql v1.5.7 // indirect
	gorm.io/gorm v1.25.11 // indirect
)
`
	filekit.FilePutContents(filePath, content)
}

func genBootInit(workspaceName, bootInitFile, gitkeepFile string) {
	content := `package boot

import (
	"` + workspaceName + `/app"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/textthree/cvgokit/filekit"
	"github.com/textthree/provider"
	"github.com/textthree/provider/clog"
	"github.com/textthree/provider/config"
)

var HttpServerPort string

func init() {
	app.Config = provider.Services.NewSingle(config.Name).(config.Service)
	HttpServerPort = app.Config.GetHttpPort()
	app.Log = provider.Services.NewSingle(clog.Name).(clog.Service)
}

func Recover(fiberApp *fiber.App) {
	if !app.Config.IsDebug() {
		fiberApp.Use(recover.New())
	}
}`
	filekit.FilePutContents(bootInitFile, content)
	filekit.DeleteFile(gitkeepFile)
}
