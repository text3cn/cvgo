package work

import (
	"cvgo/config"
	"cvgo/ins"
	"cvgo/kvs"
	"cvgo/kvs/kvsKey"
	"cvgo/tpl"
	"github.com/textthree/cvgokit/filekit"
	"github.com/textthree/provider/clog"
	"os"
	"path/filepath"
)

var workspacePath string

func CreateWork(workspaceDirName string) {
	workspacePath = filepath.Join(filekit.Getwd(), workspaceDirName)
	filekit.MkDir(workspacePath)
	createGoWorkFile()
	createGoModFile(workspaceDirName)
	CopyDirectoryStructure()
	saveInfo(workspaceDirName)
	clog.GreenPrintln("If the workspace is created successfully")
}

// create go.work file
func createGoWorkFile() {
	goWorkContent := `go ` + config.GoVersion + `

use (
	./.
)`
	// create go.work file
	file, err := os.Create(filepath.Join(workspacePath, "go.work"))
	if err != nil {
		ins.Log.Error("Create go.work file fail.", err)
		return
	}
	defer file.Close()

	// Writes content to file
	_, err = file.WriteString(goWorkContent)
	if err != nil {
		ins.Log.Error("Put content to go.work file fail.", err)
		return
	}
}

// Create go.mod file
func createGoModFile(workspaceDirName string) {
	content := `module ` + workspaceDirName + `

go ` + config.GoVersion
	content += `
require github.com/textthree/provider ` + config.CvgoProviderVersion + `

require (
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
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
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
`
	filePath := filepath.Join(workspacePath, "go.mod")
	filekit.FilePutContents(filePath, content)
}

// Copy project template
func CopyDirectoryStructure() {
	err := tpl.CopyDirFromEmbedFs(tpl.WorkDirTpl, "work", workspacePath)
	if err != nil {
		ins.Log.Error("Copy directory structure fail.", err)
	}
}

func saveInfo(workspaceDirName string) {
	kvs.Instance(workspacePath).Set(kvsKey.WorkspacePath, workspacePath)
	kvs.Instance(workspacePath).Set(kvsKey.WorkspaceName, workspaceDirName)

}
