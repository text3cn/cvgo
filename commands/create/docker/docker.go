package docker

import (
	"cvgo/paths"
	"cvgo/tpl"
	"fmt"
	"github.com/textthree/cvgokit/filekit"
	"github.com/textthree/provider/clog"
)

// 添加用 docker-compse 用于启动开发环境
// go build -o $GOPATH/bin/cvg
// cvg create-docker
func CreateDocker() {
	workPath := paths.NewWorkPath()
	cvgoPath := paths.NewCvgoPath()
	if filekit.FileExist(workPath.DockerComposeEnv()) {
		clog.RedPrintln("目标文件已存在，无法生成。", workPath.DockerComposeEnv())
		return
	}
	err := tpl.CopyFileFromEmbed(tpl.DockerComposeEnv, cvgoPath.DockerComposeEnv(), workPath.DockerComposeEnv())
	if err != nil {
		fmt.Println(err)
	}
	err = tpl.CopyDirFromEmbedFs(tpl.DockerDir, cvgoPath.DockerDir(), workPath.DockerDir())
	if err != nil {
		clog.RedPrintln(err)
	}
}
