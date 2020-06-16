package models

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func SetupModels() *gorm.DB {
	db, err := gorm.Open("mysql", "admin:xg97SmkXMGxlzztLRkTz@tcp(database-1.cdhqvttn2w7u.ap-south-1.rds.amazonaws.com:3306)/collection")

	if err != nil {
		log.Println("Error is", err)
		panic("Failed to connect to database!")
	}

	db.AutoMigrate(&Songs{})
	return db
}

