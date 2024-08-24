package mysql

import "gorm.io/gorm"

// 模型定义文档：https://gorm.io/zh_CN/docs/models.html
type BlogArticleEntity struct {
	CommonField
}

// 博客文章
func (BlogArticleEntity) TableName() string {
	return "blog_article"
}

// 添加表注释
func (this BlogArticleEntity) SetTableComment(db *gorm.DB) {
	db.Exec("ALTER TABLE blog_article COMMENT = '博客文章'")
}
