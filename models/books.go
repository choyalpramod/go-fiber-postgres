package models

import "gorm.io/gorm"

type Books struct {
	ID        uint    `json:"id" gorm:"primary_key;auto_increment"`
	Author    *string `json:"author"`
	Title     *string `json:"title"`
	Publisher *string `json:"publisher"`
}

func MigrateBooks(db *gorm.DB) error {
	return db.AutoMigrate(&Books{})
}
