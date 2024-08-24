package scripts

import (
	"cvgo/console/internal/paths"
	"github.com/textthree/cvgokit/filekit"
	"path/filepath"
)

func AddScript(script string) {
	paths.CdToProjectRootPath()
	switch script {
	case "migrate_swagger":
		migrateAndSwagger()
	}
}

// go build -o $GOPATH/bin/cvg ./console
// cvg add script migrate_swagger
func migrateAndSwagger() {
	path := paths.NewPathForProjectRoot()
	src := filepath.Join(path.ScriptsTplDir(), "migrate_swagger.go")
	dst := filepath.Join(path.AppScriptsDir(), "migrate_swagger.go")
	err := filekit.CopyFile(src, dst)
	if err == nil {
		filekit.DeleteFile(filepath.Join(path.AppScriptsDir(), ".gitkeep"))
	}
}
