package gencvgo

import (
	"cvgo/config"
	"cvgo/kvs"
	"cvgo/kvs/kvsKey"
	"cvgo/paths"
	"github.com/textthree/cvgokit/filekit"
	"github.com/textthree/cvgokit/strkit"
	"github.com/textthree/provider"
	"path/filepath"
)

var log = provider.Clog()

// 使用 cvgo http 服务
func CreateWebserverWithCvgo(modName string) {
	path := paths.NewModulePath(modName)
	moduleDir := path.ModulePath
	workspaceName, err := kvs.Instance().GetString(kvsKey.WorkspaceName)

	// 创建 index.api.go
	content := `package api

import "github.com/textthree/cvgoweb"

func Index(ctx *httpserver.Context) {
	ctx.Resp.Text("Index")
}
`
	filekit.FilePutContents(path.IndeApiGo(), content)
	filekit.DeleteFile(path.IndeApiGitkeep())

	// boot init
	genBootInit(workspaceName, path.BootInitGo(), path.BootInitGitkeep())

	// 创建路由
	content = `package routing

import (
	"` + modName + `/internal/api"
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
		authGroup.Get("/demo", api.Index)
	}
}`
	filekit.FilePutContents(path.RoutingGo(), content)
	filekit.DeleteFile(path.RoutingGitkeep())

	genAuthMiddware(path.AuthMiddlewareGo(), workspaceName)

	// 创建 main.go
	mainGoFilePath := filepath.Join(moduleDir, "main.go")
	content = `package main

import (
	_ "` + modName + `/internal/boot"
	"` + modName + `/internal/routing"
	"github.com/textthree/cvgoweb/cvgohttp"
)

func main() {
	cvgohttp.Run(routing.Routes)
}`
	filekit.FilePutContents(mainGoFilePath, content)

	// 模块 app.yaml 中添加 tokenSecret
	tokenSecret := strkit.CreateNonceStr(16)
	content = `
tokenSecret: ` + tokenSecret + `
`
	err = filekit.FileAppendContent(path.ConfigAppYaml(), content)
	if err != nil {
		log.Error(err.Error())
	}

	// go.mod
	createGoModFile(moduleDir, modName)
}

func genAuthMiddware(filepath, workspaceName string) {
	content := `package middleware

import (
	"errors"
	"github.com/spf13/cast"
	"github.com/textthree/cvgokit/cryptokit"
	"github.com/textthree/cvgokit/jsonkit"
	"github.com/textthree/cvgoweb"
	"github.com/textthree/provider"
	"github.com/textthree/provider/config"
	"` + workspaceName + `/common/dto"
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
		//clog.PinkPrintf("tokenKey=%s, token=%s, userId=%s \n", secret, token, uid)
		if token == "" || cast.ToInt64(uid) == 0 {
			ret := dto.BaseRes{
				ApiCode:    1000,
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
`
	filekit.FilePutContents(filepath, content)
}

// 创建 go.mod 文件
func createGoModFile(modDir, modName string) {
	filePath := filepath.Join(modDir, "go.mod")
	content := `module ` + modName + `

go ` + config.GoVersion
	content += `
require (
	github.com/swaggo/swag v1.16.3
	github.com/textthree/cvgoweb ` + config.CvgowebVersion + `
	github.com/textthree/provider ` + config.CvgoProviderVersion + `
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/VictoriaMetrics/fastcache v1.12.2 // indirect
	github.com/allegro/bigcache v1.2.1 // indirect
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
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/redis/go-redis/v9 v9.6.1 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.7.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.19.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/textthree/cvgokit ` + config.CvgoKitVersion + ` // indirect
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
	"github.com/textthree/provider"
	"github.com/textthree/provider/clog"
	"github.com/textthree/provider/config"
	"github.com/textthree/cvgokit/filekit"
)

func init() {
	app.Config = provider.Services.NewSingle(config.Name).(config.Service)
	app.Log = provider.Services.NewSingle(clog.Name).(clog.Service)
	clog.CyanPrintln("  Current Path: " + filekit.Getwd())
}`
	filekit.FilePutContents(bootInitFile, content)
	filekit.DeleteFile(gitkeepFile)
}
