package entity

import (
	"gorm.io/gorm"
	"text3/entity/mysql"
)

func AutoMigrate(db *gorm.DB) {
	var entities = []interface{}{}
    entities = append(entities, &mysql.BlogArticleEntity{})
	entities = append(entities, &mysql.BlogArticleCateEntity{})
	entities = append(entities, &mysql.UserEntity{})

	db.AutoMigrate(
		entities...,
	)
	addTableComment(db)
}

func addTableComment(db *gorm.DB) {
    mysql.BlogArticleEntity{}.SetTableComment(db)
	mysql.BlogArticleCateEntity{}.SetTableComment(db)
	mysql.UserEntity{}.SetTableComment(db)

}
