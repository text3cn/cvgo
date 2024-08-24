package mysql

import "gorm.io/gorm"

// 模型定义文档：https://gorm.io/zh_CN/docs/models.html
type BlogArticleCateEntity struct {
	CommonField
	Pid      int    `gorm:"index;default:0"`
	UserId   int    `gorm:"index;default:0"`
	Title    string `gorm:"index; comment:分类名称"`
	Expand   int8   `gorm:"comment:是否展开子分类;default:1"`
	Selected int8   `gorm:"comment:是否选中状态;default:0"`
	Sort     int    `gorm:"index;comment:排序;default:0"`
}

// 博文分类
func (BlogArticleCateEntity) TableName() string {
	return "blog_article_cate"
}

// 添加表注释
func (this BlogArticleCateEntity) SetTableComment(db *gorm.DB) {
	db.Exec("ALTER TABLE blog_article_cate COMMENT = '博文分类'")
}
