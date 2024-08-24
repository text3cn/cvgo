package entityregistry

import (
	"reflect"
)

// 全局类型注册表，用于 CURD 代码生成
var EntityRegistry = make(map[string]reflect.Type)

func init() {

}
