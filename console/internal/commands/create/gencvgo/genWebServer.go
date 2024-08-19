package gencvgo

import (
	"cvgo/kit/filekit"
	"cvgo/kit/strkit"
	"cvgo/provider"
	"path/filepath"
)

var log = provider.Clog()

// 使用 cvgo http 服务
func CreateWebserverWithCvgo(moduleDir, modName string) {
	// 创建 index.api.go
	apiFilePath := filepath.Join(moduleDir, "internal", "api", "index.api.go")
	content := `package api

import "cvgo/provider/httpserver"

func Index(ctx *httpserver.Context) {
	ctx.Resp.Text("Index")
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
	"cvgo/provider/httpserver"
	"cvgo/provider/httpserver/middleware"
)

// 路由定义
func Routes(engine *httpserver.Engine) {
	engine.Get("/", api.Index)

	// 需要验证 token 的路由分组
	authGroup := engine.Prefix("/api").UseMiddleware(middleware.Auth())
	{
		authGroup.Get("/demo", api.Index)
	}
}`
	filekit.FilePutContents(routingFilePath, content)
	err = filekit.DeleteFile(filepath.Join(moduleDir, "internal", "routing", ".gitkeep"))
	if err != nil {
		log.Error(err.Error())
	}
	// 创建 main.go
	mainGoFilePath := filepath.Join(moduleDir, "main.go")
	content = `package main

import (
	_ "` + modName + `/internal/boot"
	"` + modName + `/internal/routing"
	"cvgo/provider/httpserver/cvgohttp"
)

func main() {
	cvgohttp.Run(routing.Routes)
}`
	filekit.FilePutContents(mainGoFilePath, content)

	// 模块 app.yaml 中添加 tokenSecret
	appYamlFile := filepath.Join(moduleDir, "internal", "config", "app.yaml")
	tokenSecret := strkit.CreateNonceStr(16)
	content = `
tokenSecret: ` + tokenSecret + `
`
	err = filekit.FileAppendContent(appYamlFile, content)
	if err != nil {
		log.Error(err.Error())
	}
}
