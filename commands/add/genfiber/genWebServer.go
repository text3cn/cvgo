package genfiber

import (
	"cvgo/console/internal/paths"
	"github.com/textthree/cvgokit/filekit"
	"github.com/textthree/provider"
	"path/filepath"
)

var log = provider.Clog()

// 使用 fiber http 服务
func CreateWebserver(moduleDir, modName string, supportSwagger bool) {
	// 创建 index.api.go
	apiFilePath := filepath.Join(moduleDir, "internal", "api", "index.api.go")
	content := `package api

import "github.com/gofiber/fiber/v2"

func Index(ctx *fiber.Ctx) error {
    return ctx.SendString("Index")
}
`
	filekit.FilePutContents(apiFilePath, content)
	err := filekit.DeleteFile(filepath.Join(moduleDir, "internal", "api", ".gitkeep"))
	if err != nil {
		log.Error(err.Error())
	}

	// 创建路由
	routingFilePath := filepath.Join(moduleDir, "internal", "routing", "routing.go")
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
	filekit.FilePutContents(routingFilePath, content)
	err = filekit.DeleteFile(filepath.Join(moduleDir, "internal", "routing", ".gitkeep"))
	if err != nil {
		log.Error(err.Error())
	}

	// 鉴权中间件
	src := filepath.Join(paths.FiberTplForRoot, "middleware", "auth.go")
	dest := filepath.Join(moduleDir, "internal", "middleware", "auth.go")
	filekit.DeleteFile(filepath.Join(moduleDir, "internal", "middleware", ".gitkeep"))
	filekit.CopyFile(src, dest)

	// boot -> init.go 中添加 Recover 中间件
	initFile := filepath.Join(moduleDir, "internal", "boot", "init.go")
	content = `    "github.com/gofiber/fiber/v2/middleware/recover"
    "github.com/gofiber/fiber/v2"`
	filekit.AddContentUnderLine(initFile, "import (", content)
	content = `
func Recover(fiberApp *fiber.App) {
	if !app.Config.IsDebug() {
		fiberApp.Use(recover.New())
	}
}
`
	filekit.FileAppendContent(initFile, content)

	// 创建 main.go
	swagger := ""
	if supportSwagger {
		src = filepath.Join(paths.FiberTplForRoot, "boot", "swagger.go")
		dest = filepath.Join(moduleDir, "internal", "boot", "swagger.go")
		filekit.CopyFile(src, dest)
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
}
