package common

import (
	"github.com/textthree/cvgokit/filekit"
	"github.com/textthree/provider/clog"
	"strings"
)

// 在一个 go 文件中判断如果包没有导入则导报
func ImportPackageIfNotImport(filePath, packagePath string) {
	fileContent, _ := filekit.FileGetContents(filePath)
	if !strings.Contains(fileContent, `"`+packagePath+`"`) {
		content := `    "` + packagePath + `"`
		err := filekit.AddContentUnderLine(filePath, "import (", content)
		if err != nil {
			clog.RedPrintln(err)
		}
	}
}
