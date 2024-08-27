package gitlabci

import (
	"cvgo/kvs"
	"cvgo/paths"
	"cvgo/tpl"
	"fmt"
	"github.com/textthree/cvgokit/filekit"
	"github.com/textthree/provider/clog"
	"path/filepath"
)

// go build -o $GOPATH/bin/cvg
// cvg create-gitlab-ci
func CreateGitlabCiYml() {
	workdir := kvs.Instance().GetWorkspacePath()
	cvgoPath := paths.NewCvgoPath()
	fmt.Printf("workdir: %s\n", workdir+cvgoPath.GitlabCI())
	dst := filepath.Join(workdir, cvgoPath.GitlabCI())
	if filekit.FileExist(dst) {
		clog.RedPrintln("目标文件已存在，无法生成。", dst)
		return
	}
	err := tpl.CopyFileFromEmbed(tpl.GitlabCI, cvgoPath.GitlabCI(), dst)
	if err != nil {
		fmt.Println(err)
	}
}
