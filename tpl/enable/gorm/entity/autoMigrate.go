package entity

import (
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	var entities = []interface{}{}

	db.AutoMigrate(
		entities...,
	)
	addTableComment(db)
}

func addTableComment(db *gorm.DB) {

}
