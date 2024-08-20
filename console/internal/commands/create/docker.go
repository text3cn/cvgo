package create

import (
	"cvgo/console/internal/paths"
	"cvgo/kit/filekit"
	"fmt"
)

// 添加用 docker-compse 用于启动开发环境
// go build -o $GOPATH/bin/cvg ./console && cvg create dockerEnv
func createDockerEnv() {
	paths.NewPathForApp()
	path := paths.NewPathForAnywhere()
	err := filekit.CopyFile(path.DockerEnvTpl(), path.AppDockerEnv())
	if err != nil {
		fmt.Println(err)
	}
	filekit.CopyFiles(path.DockerDirTpl(), path.AppDockerDir())
}
