package create

import (
	"bufio"
	"bytes"
	"cvgo/console/console"
	"cvgo/kit/filekit"
	"cvgo/provider/clog"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var pwd string

// 创建模块
func createModule(modName, webserver string, swagger, force bool) {
	var err error
	pwd, err = os.Getwd()
	if err != nil {
		log.Error("os.Getwd() 失败", err.Error())
		return
	}
	// 检查是否是在根目录
	goWorkSumFile := filepath.Join(pwd, "go.work.sum")
	if !filekit.FileExist(goWorkSumFile) {
		log.Error("请在根目录运行此命令")
		return
	}
	// 检查是否存在 go.work
	goWorkFile := filepath.Join(filekit.Getwd(), "go.work")
	if !filekit.FileExist(goWorkFile) {
		createGoWorkFile()
	}
	// 检查是否存在 app 目录
	var isFirstInit bool
	if !filekit.DirExists(filepath.Join(pwd, "app")) {
		isFirstInit = true
	}
	// 创建模块
	err = initModule(modName, isFirstInit, force)
	if err != nil {
		if err.Error() == "cancel" {
			log.Info("你取消了操作")
			return
		}
		log.Error(err.Error())
		return
	}
	createCodeFile(modName, webserver, swagger) // 创建代码文件
	// go.work 增加 module
	appendModuleToGoWorkFile(modName)
}

// 创建 go.work 文件
func createGoWorkFile() {
	goWorkContent := `go ` + console.GoVersion + `

use (
	./.
)`
	// 创建 go.work 文件
	file, err := os.Create("go.work")
	if err != nil {
		log.Error("创建 go.work 失败", err)
		return
	}
	defer file.Close()

	// 写入内容到文件
	_, err = file.WriteString(goWorkContent)
	if err != nil {
		log.Error("往 go.work 写入内容失败", err)
		return
	}
}

// 创建模块
// force 是否强制创建，如果目标模已经存在则先删除再创建
func initModule(modName string, isFirstInit, force bool) error {
	dir := filepath.Join(pwd, "app", "modules", modName)
	if filekit.DirExists(dir) {
		if force {
			if goOn := deleteBeforeCreate(modName); goOn == false {
				return errors.New("cancel")
			}
		} else {
			return errors.New("模块 " + modName + " 已经存在")
		}
	}
	filekit.MkDir(dir)
	createGoModFile(modName, filepath.Join(dir, "go.mod"))
	src := filepath.Join(pwd, "console", "commands", "create", "sampletpl", "mod")
	filekit.CopyFiles(src, filepath.Join(dir)) // 拷贝 module 模板
	// 拷贝 app 模板
	if isFirstInit {
		src = filepath.Join(pwd, "console", "commands", "create", "sampletpl", "app")
		filekit.CopyFiles(src, filepath.Join(pwd, "app"))
	}
	return nil
}

func deleteBeforeCreate(moduleName string) (delete bool) {
	cur, _ := os.Getwd()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(clog.YellowSprintf("⚠️ 模块" + moduleName + "已经存在，是否将其删除后重新创建？(yes/no) [default:" + clog.CyanSprintf("yes", clog.ColorYellow) + "]:"))
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "yes" {
		delete = true
	} else if input == "no" {
		delete = false
	} else {
		delete = true
	}
	if delete {
		filekit.DeleteDirOrFile(filepath.Join(cur, "app", "modules", moduleName))
		err := deleteModuleInGoWork(moduleName, filepath.Join(cur, "go.work"))
		if err != nil {
			log.Error(err.Error())
			delete = false
		}
	}
	return
}

// 创建 go.mod 文件
func createGoModFile(modName, filePath string) {
	content := `module ` + modName + `

go ` + console.GoVersion
	filekit.FilePutContents(filePath, content)
}

// go.work 中追加 module
func appendModuleToGoWorkFile(modName string) {
	// 读取现有的 go.work 文件内容
	filePath := "go.work"
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// 读取文件内容并准备修改
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// 查找并在 `./.` 之后插入新的路径
	for i, line := range lines {
		if strings.TrimSpace(line) == "./." {
			// 在 `./.` 之后插入 `./app/modules/modName`
			lines = append(lines[:i+1], append([]string{"\t./app/modules/" + modName}, lines[i+1:]...)...)
			break
		}
	}

	// 将修改后的内容写回到文件
	outputFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}

	// 确保所有缓冲内容写入文件
	if err := writer.Flush(); err != nil {
		fmt.Println("Error flushing writer:", err)
		return
	}
}

// 从 go.work 中删除指定模块
func deleteModuleInGoWork(moduleName, filePath string) error {
	targetLine := "./app/modules/" + moduleName
	// 打开文件进行读取
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	var outputLines []string
	scanner := bufio.NewScanner(file)

	// 遍历文件的每一行
	for scanner.Scan() {
		line := scanner.Text()
		// 如果当前行包含需要删除的导入路径，则跳过
		if strings.TrimSpace(line) == targetLine {
			continue
		}
		outputLines = append(outputLines, line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取文件时出错: %v", err)
	}

	// 将修改后的内容写回文件
	if err := os.WriteFile(filePath, []byte(strings.Join(outputLines, "\n")+"\n"), 0644); err != nil {
		return fmt.Errorf("写入文件时出错: %v", err)
	}

	return nil
}

// 创建代码文件
func createCodeFile(modName, webserver string, supportSwagger bool) {
	moduleDir := filepath.Join(pwd, "app", "modules", modName)

	if webserver == "" {
		createEmptyMainGoFile(moduleDir, modName)
	} else {
		if webserver == "cvgo" {
			createWebserverWithCvgo(moduleDir, modName)
		}
		createAppYaml(moduleDir, supportSwagger)
	}

	// swagger 支持
	if supportSwagger {
		// 在 main 函数上添加 swagger 注释
		content := `
// @title ` + modName + ` API Doc
// @description 使用 [swaggo](https://github.com/swaggo/swag/blob/master/README_zh-CN.md) 编写。
// @host localhost:` + console.DefaultHttpPort + `
// @BasePath /
`
		fmt.Println(cfg)
		moddir := filepath.Join(pwd, "app", "modules", modName)
		filekit.AddContentAboveLine(filepath.Join(moddir, "main.go"), "func main() {", content)

		// swag init
		err := os.Chdir(moddir)
		if err != nil {
			log.Error("进入", moddir, "失败", err.Error())
			return
		}
		cmd := exec.Command("swag", "init", "--parseDependency", "--propertyStrategy", "pascalcase")
		var out bytes.Buffer
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			log.Error("添加 Swagger 失败:", err.Error())
			log.Info("请确保你已安装 swaggo 然后重新创建 module。", "https://github.com/swaggo/swag")
			os.Chdir(pwd)
			return
		}
		os.Chdir(pwd)

	}
}

// 创建 app.yaml
func createAppYaml(moduleDir string, supportSwagger bool) {
	filePath := filepath.Join(moduleDir, "internal", "config", "app.yaml")
	content := `server:
  http-port: ` + console.DefaultHttpPort + ` # http 服务监听的端口
`
	if supportSwagger {
		content += `
swagger:
  filepath: ./docs/swagger.json
`
	}
	filekit.FilePutContents(filePath, content)
	filekit.DeleteFile(filepath.Join(moduleDir, "internal", "config", ".gitkeep"))
}

// 空 main.go
func createEmptyMainGoFile(moduleDir, modName string) {
	filePath := filepath.Join(moduleDir, "main.go")
	content := `package main

import (
	_ "` + modName + `/internal/boot"
)

func main() {

}
`
	filekit.FilePutContents(filePath, content)
}

// 使用 cvgohttp 服务
func createWebserverWithCvgo(moduleDir, modName string) {
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
)

// 路由定义
func Router(engine *httpserver.Engine) {
	engine.Get("/", api.Index)
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
	cvgohttp.Run(routing.Router)
}`
	filekit.FilePutContents(mainGoFilePath, content)
}
