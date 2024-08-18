package types

import "cvgo/provider/core"

type BeforStartCallback func(c core.Container)

// clog 配置
type Clog struct {
	Level string
}

// 阿里 oss 配置
type AliOss struct {
	OSS_ID string
}
