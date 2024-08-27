package gencvgo

import (
	"cvgo/commands/codegen/common"
	"cvgo/kvs"
	"cvgo/kvs/kvsKey"
	"cvgo/paths"
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
func GenApi(method, requestPath string, supportSwagger, cursorPaging bool, cvgflag string, svc ...string) {
	var tableName, svcName, svcFuncName, curdlType string
	if len(svc) > 0 {
		tableName = svc[0]
		svcName = svc[1]
		svcFuncName = svc[2]
		curdlType = svc[3]
	}
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
	if arrkit.InArray(requestPath+"["+method+"]", oldRoutes) {
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
	funcName := methodToCurd(method, curdlType) + strkit.Ucfirst(pathArr[1]) + strkit.Ucfirst(pathArr[2])

	routePath := pathArr[1] + "/" + strkit.CamelToKebabCase(pathArr[2])
	content := `` + addSwagger(supportSwagger, funcName, routePath, method) + `
func ` + funcName + ` (ctx *httpserver.Context) {
	` + addApiFuncCode(funcName, tableName, svcName, svcFuncName, curdlType, cursorPaging) + `
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
	cursor := ""
	if cursorPaging {
		cursor = "    Cursor int64"
	}
	content = `
type ` + funcName + `Req struct {
` + cursor + `
}
	
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
		content = space + cvgflag + `.` + strkit.Ucfirst(method) + `("` + routePath + `", api.` + funcName + `)`
		err = filekit.AddContentUnderLine(path.RoutingGo(), "// cvgflag="+cvgflag, content)
		if err != nil {
			clog.RedPrintln(err)
		}
	}

	// 生 apidebug
	fileName := pathArr[2] + strkit.Ucfirst(method) + ".html"
	apiDebugHtmlFile := filepath.Join(path.ModuleApiDebugDir(), pathArr[1], fileName)
	common.GenApidebug(apiDebugHtmlFile, requestPath, method, cursorPaging)

	// 完成
	kv.Set(modName+".routes", append(oldRoutes, requestPath+"["+method+"]"))
	clog.GreenPrintln("创建 API 成功：", pathArr[1]+".api.go", "->", funcName+"()")
}

// qpi 基础代码
func addApiFuncCode(funcName, tableName, svcName, svcFuncName, curdlType string, cursorPaging bool) string {
	content := `req := dto.` + funcName + `Req{}
	res := dto.` + funcName + `Res{}
	if err := ctx.Req.JsonScan(&req); err != nil {
		ctx.Resp.Json(cvgerr.ParseRequestParamsError())
		return
	}`

	if tableName != "" {
		switch curdlType {
		case "c":
			content += `
	if err := service.` + strkit.Ucfirst(svcName) + `Svc(ctx).` + svcFuncName + `(); err != nil {
		res.ApiCode = cvgerr.Fail
		res.ApiMessage = err.Error()
	}`

		case "u":
			content += `
	if err := service.` + strkit.Ucfirst(svcName) + `Svc(ctx).` + svcFuncName + `(); err != nil {
		res.ApiCode = cvgerr.Fail
		res.ApiMessage = err.Error()
	}`

		case "r":
			content += `
	service.` + strkit.Ucfirst(svcName) + `Svc(ctx).` + svcFuncName + `()`

		case "d":
			content += `
	if err := service.` + strkit.Ucfirst(svcName) + `Svc(ctx).` + svcFuncName + `(); err != nil {
		res.ApiCode = cvgerr.Fail
		res.ApiMessage = err.Error()
	}`

		case "l":
			cursor := ""
			if cursorPaging {
				cursor = "req.Cursor"
			}
			content += `
	service.` + strkit.Ucfirst(svcName) + `Svc(ctx).` + svcFuncName + `(` + cursor + `)`

		}
	}

	content += `
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

// 根据 curdl 转换成 Create / Update / Get/ Delete / List
func methodToCurd(method, curdlType string) string {
	if curdlType == "" {
		return strkit.Ucfirst(method)
	}
	ret := ""
	switch curdlType {
	case "c":
		ret = "Create"
	case "u":
		ret = "Update"
	case "r":
		ret = "Get"
	case "d":
		ret = "Delete"
	case "l":
		ret = "List"
	}
	return ret
}
