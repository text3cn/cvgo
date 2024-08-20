package entity

import (
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	var entitys = []interface{}{}

	db.AutoMigrate(
		entitys...,
	)
}
