package module

import (
	"bufio"
	"bytes"
	"cvgo/commands/codegen/gencvgo"
	"cvgo/commands/codegen/genfiber"
	"cvgo/ins"
	"cvgo/kvs"
	"cvgo/kvs/kvsKey"
	"cvgo/tpl"
	"fmt"
	"github.com/spf13/cast"
	"github.com/textthree/cvgokit/filekit"
	"github.com/textthree/provider/clog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var kvsIns kvs.KvStorage
var modDir string
var httpPort string
var workDir string

// go build -o $GOPATH/bin/cvg
// cvg create-module chord --webserver=cvgo --force --swagger
func CreateModule(modName, webserver string, swagger, force bool) {
	kvsIns = kvs.Instance()
	workDir = kvsIns.GetRootPath()
	var err error

	// Allocate port
	portInt, _ := kvsIns.GetAllocatedPort()
	httpPort = cast.ToString(portInt + 1)
	kvsIns.Set(kvsKey.AllocatedPort, portInt+1)

	// Create
	modDir = filepath.Join(workDir, "app", modName)
	if filekit.DirExists(modDir) {
		if force {
			if goOn := deleteBeforeCreate(workDir, modDir, modName); goOn == false {
				if err != nil {
					if err.Error() == "cancel" {
						clog.CyanPrintln("You canceled the operation")
						return
					}
					return
				}
				return
			}
		} else {
			clog.RedPrintln("Module " + modName + " already exists")
			return
		}
	}

	filekit.MkDir(modDir)

	// Copy module template
	err = tpl.CopyDirFromEmbedFs(tpl.ModuleTpl, "module", modDir)
	if err != nil {
		panic(err)
	}

	createCodeFile(modName, webserver, swagger)
	appendModuleToGoWorkFile(modName)

	// Save info
	kvsIns.Set(kvsKey.ModuleWebFramework(modName), webserver)
	kvsIns.Set(kvsKey.ModuleHttpPort(modName), httpPort)
	kvsIns.Set(kvsKey.ModuleSwaggerEnable(modName), swagger)
	clog.GreenPrintln("Create module", modName, "successfully.")
}

func deleteBeforeCreate(workDir, modDir, moduleName string) (delete bool) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(clog.YellowSprintf("⚠️ Module" + moduleName + " already exists. Do you want to delete it and create it again? (yes/no) [default:" + clog.CyanSprintf("yes", clog.ColorYellow) + "]:"))
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
		err := filekit.DeleteDirOrFile(modDir)
		if err != nil {
			clog.RedPrintln(err)
		}
		err = deleteModuleInGoWork(moduleName, filepath.Join(workDir, "go.work"))
		if err != nil {
			ins.Log.Error(err.Error())
			delete = false
		}
	}
	return
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
			lines = append(lines[:i+1], append([]string{"\t./app/" + modName}, lines[i+1:]...)...)
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
	targetLine := "./app/" + moduleName
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

func createCodeFile(modName, webserver string, supportSwagger bool) {
	if webserver == "" {
		createEmptyMainGoFile(modDir, modName)
	} else {
		createAppYaml(modDir, supportSwagger)
		if webserver == "cvgo" {
			gencvgo.CreateWebserverWithCvgo(modName)
		} else if webserver == "fiber" {
			genfiber.CreateWebserver(modName, supportSwagger)
		}
	}

	// swagger
	if supportSwagger {
		// Add a swagger comment to the main function
		content := `
// @title ` + modName + ` API Doc
// @description 使用 [swaggo](https://github.com/swaggo/swag/blob/master/README_zh-CN.md) 编写。
// @host localhost:` + httpPort + `
// @BasePath /
`
		filekit.AddContentAboveLine(filepath.Join(modDir, "main.go"), "func main() {", content)

		// swag init
		err := os.Chdir(modDir)
		if err != nil {
			panic(err.Error())
			return
		}
		cmd := exec.Command("swag", "init", "--parseDependency", "--propertyStrategy", "pascalcase")
		var out bytes.Buffer
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			clog.RedPrintln("Make sure you have swaggo installed and then re-create the module.", "https://github.com/swaggo/swag")
			panic(err.Error())
			os.Chdir(workDir)
			return
		}
		os.Chdir(workDir)

	}
}

// app.yaml
func createAppYaml(moduleDir string, supportSwagger bool) {
	filePath := filepath.Join(moduleDir, "internal", "config", "app.yaml")
	content := `server:
  http-port: ` + httpPort + `  
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
