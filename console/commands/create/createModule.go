package create

import (
	"bufio"
	"cvgo/kit/filekit"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var goVersion = "1.22.4"
var pwd string

// 创建模块
func createModule(modName string) {
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
	err = initModule(modName, isFirstInit)
	if err != nil {
		log.Error(err.Error())
		return
	}
	// go.work 增加 module
	appendModuleToGoWorkFile(modName)
}

// 创建 go.work 文件
func createGoWorkFile() {
	goWorkContent := `go ` + goVersion + `

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
func initModule(modName string, isFirstInit bool) error {
	dir := filepath.Join(pwd, "app", "modules", modName)
	if filekit.DirExists(dir) {
		return errors.New("模块 " + modName + " 已经存在")
	}
	filekit.MkDir(dir)
	createGoModFile(modName, filepath.Join(dir, "go.mod"))
	src := filepath.Join(pwd, "console", "commands", "create", "sampletpl", "mod")
	filekit.CopyFiles(src, filepath.Join(dir)) // 拷贝 module 模板
	createMainGoFile(modName, filepath.Join(dir, "main.go"))
	// 拷贝 app 模板
	if isFirstInit {
		src = filepath.Join(pwd, "console", "commands", "create", "sampletpl", "app")
		filekit.CopyFiles(src, filepath.Join(pwd, "app"))
	}
	return nil
}

// 创建 go.mod 文件
func createGoModFile(modName, filePath string) {
	content := `module ` + modName + `

go ` + goVersion
	filekit.FilePutContents(filePath, content)
}

// 创建 main.go 文件
func createMainGoFile(modName, filePath string) {
	content := `package main

import (
	_ "` + modName + `/internal/boot"
)

func main() {

}
`
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
