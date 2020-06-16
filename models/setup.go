package models

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func SetupModels() *gorm.DB {
	db, err := gorm.Open("mysql", "")

	if err != nil {
		log.Println("Error is", err)
		panic("Failed to connect to database!")
	}

	db.AutoMigrate(&Songs{})
	return db
}
