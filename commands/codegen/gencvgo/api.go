package gencvgo

import (
	"cvgo/commands/codegen/common"
	"cvgo/kvs"
	"cvgo/kvs/kvsKey"
	"cvgo/paths"
	"fmt"
	"github.com/textthree/cvgokit/arrkit"
	"github.com/textthree/cvgokit/filekit"
	"github.com/textthree/cvgokit/gokit"
	"github.com/textthree/cvgokit/strkit"
	"github.com/textthree/provider/clog"
	"path/filepath"
	"strings"
)

// 模块目录下执行：
// go build -o $GOPATH/bin/cvg
// cvg create-api get api/user/info
func GenApi(method, requestPath string, supportSwagger bool, cvgflag string) {
	path := paths.NewModulePath()
	paths.CheckRunAtModuleRoot()
	modName, _ := gokit.GetModuleName()
	kv := kvs.Instance()
	pathArr := strkit.Explode("/", requestPath)
	if len(pathArr) != 3 || pathArr[2] == "" {
		clog.RedPrintln("请求路径格式错误，路径为三段。示例：api/user/info")
		return
	}
	oldRoutes, _ := kv.GetStringSlice(kvsKey.ModuleRoute(modName))
	if arrkit.InArray(requestPath, oldRoutes) {
		clog.RedPrintln("API", requestPath, "已存在")
		return
	}

	// 创建 api
	workspaceName := kv.GetWorkspaceName()
	apiFile := filepath.Join(path.ModuleApiDir(), pathArr[1]+".api.go")
	err := filekit.CreatePath(apiFile)
	if err == nil {
		content := `package api
	
import (
	"chord/internal/dto"
	"` + workspaceName + `/cvgerr"
	"github.com/textthree/cvgoweb"
)
`
		filekit.FilePutContents(apiFile, content)
	}
	funcName := strkit.Ucfirst(method) + strkit.Ucfirst(pathArr[1]) + strkit.Ucfirst(pathArr[2])
	routePath := pathArr[1] + "/" + pathArr[2]
	content := `` + addSwagger(supportSwagger, funcName, requestPath, method) + `
func ` + funcName + ` (ctx *httpserver.Context) {
	` + addApiFuncCode(method) + `
}
`
	filekit.FileAppendContent(apiFile, content)

	// 创建 dto
	dtoFile := filepath.Join(path.ModuleDtoDir(), pathArr[1]+".dto.go")
	err = filekit.CreatePath(dtoFile)
	if err == nil {
		filekit.DeleteFile(filepath.Join(path.ModuleDtoDir(), ".gitkeep"))
		content = `package dto
	
import (
	"` + workspaceName + `/common"
)
`
		content += "/*\n\ttype Example struct {\n\t    Field  int   `json:\"field\"  swaggertype:\"integer\" validate:\"required\" description:\"描述\" example:\"示例\"`\n\t    Field2 []int `json:\"field2\" swaggertype:\"array,number\"` // 字段注释\n\t}\n*/\n"
		filekit.FilePutContents(dtoFile, content)
	}
	content = `
type ` + funcName + `Req struct {}
	
type ` + funcName + `Res struct {
	common.BaseRes
}
`
	filekit.FileAppendContent(dtoFile, content)

	// 添加路由
	if cvgflag != "" {
		space := `    `
		if pathArr[0] != "root" {
			space = `        `
		}
		fmt.Println("// cvgflag=" + cvgflag)
		content = space + cvgflag + `.` + strkit.Ucfirst(method) + `("` + routePath + `", api.` + funcName + `)`
		err = filekit.AddContentUnderLine(path.RoutingGo(), "// cvgflag="+cvgflag, content)
		if err != nil {
			clog.RedPrintln(err)
		}
	}

	// 生 apidebug
	apiDebugHtmlFile := filepath.Join(path.ModuleApiDebugDir(), pathArr[1], pathArr[2]+".html")
	common.GenApidebug(apiDebugHtmlFile, requestPath, method)

	// 完成
	kv.Set(modName+".routes", append(oldRoutes, requestPath))
	clog.GreenPrintln("生成 API 成功")
}

// qpi 基础代码
func addApiFuncCode(method string) string {
	content := `req := dto.GetUserInfoReq{}
	res := dto.GetUserInfoRes{}
	if err := ctx.Req.JsonScan(&req); err != nil {
		ctx.Resp.Json(cvgerr.ParseRequestParamsError())
		return
	}

	ctx.Resp.Json(res)`
	return content
}

// swagger
func addSwagger(supportSwagger bool, funcName, requestPath, method string) string {
	pathArr := strings.Split(requestPath, "/")
	if supportSwagger {
		content := `
// @Summary    ` + funcName + `
// @Tags        ` + pathArr[1] + `
// @Router      /` + requestPath + ` [` + method + `]
// @Accept      json
// @Produce     json
// @Param       request body dto.` + funcName + `Req true " "
// @Success 	200 {object} dto.` + funcName + `Res`
		return content
	}
	return ""
}
