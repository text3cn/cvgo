package app

import (
	"github.com/textthree/provider/clog"
	"github.com/textthree/provider/config"
)

// 这些变量的值是在 boot -> init.go 中进行初始化

var Log clog.Service

var Config config.Service
