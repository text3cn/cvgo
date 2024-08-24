package entityregistry

import (
	"reflect"
	"text3/entity/mysql"
)

// 全局类型注册表，用于 CURD 代码生成
var EntityRegistry = make(map[string]reflect.Type)

func init() {
    EntityRegistry["blog_article"] =  reflect.TypeOf(mysql.BlogArticleEntity{})
    EntityRegistry["blog_article_cate"] =  reflect.TypeOf(mysql.BlogArticleCateEntity{})
	EntityRegistry["user"] = reflect.TypeOf(mysql.UserEntity{})

}
