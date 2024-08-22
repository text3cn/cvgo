package gencvgo

import (
	"cvgo/console/internal/commands/add/addcommon"
	"cvgo/console/internal/console"
	"cvgo/console/internal/paths"
	"cvgo/kit/arrkit"
	"cvgo/kit/filekit"
	"cvgo/kit/gokit"
	"cvgo/kit/strkit"
	"cvgo/provider/clog"
	"path/filepath"
	"strings"
)

// 模块目录下执行：
// cd ../../../ && go build -o $GOPATH/bin/cvg ./console && cd app/modules/chord && cvg add svc api/user/info
func GenApi(method, requestPath string, supportSwagger bool, cvgflag string) {
	path := paths.NewPathForModule()
	modName, _ := gokit.GetModuleName()
	kv := console.NewKvStorage(filekit.GetParentDir(3))
	pathArr := strkit.Explode("/", requestPath)
	if len(pathArr) != 3 || pathArr[2] == "" {
		clog.RedPrintln("请求路径格式错误，路径为三段。示例：api/user/info")
		return
	}
	oldRoutes, _ := kv.GetStringSlice(modName + ".routes")
	if arrkit.InArray(requestPath, oldRoutes) {
		clog.RedPrintln("API", requestPath, "已存在")
		return
	}

	// 创建 api
	apiFile := filepath.Join(path.ModuleApiDir(), pathArr[1]+".api.go")
	err := filekit.CreatePath(apiFile)
	if err == nil {
		content := `package api

import (
	"chord/internal/dto"
	"cvgo/cvgerr"
	"cvgo/provider/httpserver"
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
	"cvgo/app/common/dto"
)
`
		content += "/*\n\ttype Example struct {\n\t    Field  int   `json:\"field\"  swaggertype:\"integer\" validate:\"required\" description:\"描述\" example:\"示例\"`\n\t    Field2 []int `json:\"field2\" swaggertype:\"array,number\"` // 字段注释\n\t}\n*/\n"
		filekit.FilePutContents(dtoFile, content)
	}
	content = `
type ` + funcName + `Req struct {}

type ` + funcName + `Res struct {
	dto.BaseRes
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
		err = filekit.AddContentUnderLine(path.ModuleRoutingFile(), "// cvgflag="+cvgflag, content)
		if err != nil {
			clog.RedPrintln(err)
		}
	}

	// 生 apidebug
	apiDebugHtmlFile := filepath.Join(path.ModuleApiDebugDir(), pathArr[1], pathArr[2]+".html")
	addcommon.GenApidebug(apiDebugHtmlFile, requestPath, method)

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
// @Router      ` + requestPath + ` [` + method + `]
// @Accept      json
// @Produce     json
// @Param       request body dto.` + funcName + `Req true " "
// @Success 	200 {object} dto.` + funcName + `Res`
		return content
	}
	return ""
}
