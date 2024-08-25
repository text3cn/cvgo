package app

import "os"

// 生产环境添加 ENV 环境变量：
// vim ~/.bashrc
// 文件末尾加入：export ENV=production
// 生效：source ~/.bashrc
func IsDevelop() bool {
	return os.Getenv("ENV") != "production"
}

func Env() string {
	return os.Getenv("ENV")
}
