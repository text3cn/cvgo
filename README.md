## 简介
Cvgo 是一个 Golang 后端项目开发脚手架，提供了一套基础模板与工具包用于快速搭建项目，支持代码生成、自动编译、工具包、常用函数等，
你可以在此基础上使用任何开发框架进行业务开发。

## 第三方库
- [GORM](https://gorm.io/index.html)
- [Go Redis](https://redis.uptrace.dev)
- [Cast](https://github.com/spf13/cast)
- [Viper](https://github.com/spf13/viper)
- [Cobra](https://github.com/spf13/cobra)

## 文档
开发文档：[http://cvgo.text3.cn](http://cvgo.text3.cn)
 
## 安装
### Mac
```shell
go install github.com/textthree/cvg@latest
```
将 `GOPATH/bin` 加入环境变量
```shell
echo 'export PATH="$(go env GOPATH)/bin:$PATH"' >> ~/.bash_profile
```
使环境变量生效
```shell
source ~/.bash_profile
```
然后就可以使用 cvg 命令来搭建项目了。

## 详细文档
http://cvgo.text3.cn