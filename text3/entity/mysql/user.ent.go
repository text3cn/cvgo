package mysql

import (
	"gorm.io/gorm"
	"time"
)

// 模型定义文档：https://gorm.io/zh_CN/docs/models.html
type UserEntity struct {
	CommonField
	Password   string
	Nickname   string
	Avatar     string
	AccessTime time.Time
}

// 用户表
func (UserEntity) TableName() string {
	return "user"
}

// 添加表注释
func (this UserEntity) SetTableComment(db *gorm.DB) {
	db.Exec("ALTER TABLE user COMMENT = '用户表'")
}
