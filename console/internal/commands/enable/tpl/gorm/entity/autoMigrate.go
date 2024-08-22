package entity

import (
	"gorm.io/gorm"
	"reflect"
)

// 全局类型注册表，用于 CURD 代码生成
var EntityRegistry = make(map[string]reflect.Type)

func init() {

}

func AutoMigrate(db *gorm.DB) {
	var entitys = []interface{}{}

	db.AutoMigrate(
		entitys...,
	)
	addTableComment(db)
}

func addTableComment(db *gorm.DB) {

}
