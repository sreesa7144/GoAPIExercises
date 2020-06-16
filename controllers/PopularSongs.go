package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sreesa7144/GoAPIExercises/models"
)

func PopularSongs(c *gin.Context) {
	log.Println("Coming to Popular")
	db := c.MustGet("db").(*gorm.DB)
	song := &models.Songs{}
	rows, err := db.Model(&song).Order("Popularity desc").Rows()
	var songs []models.Songs
	if err != nil {
		log.Println(err)
	} else {
		for rows.Next() {
			var song models.Songs
			db.ScanRows(rows, &song)
			songs = append(songs, song)
		}

	}

	c.JSON(http.StatusOK, gin.H{"data": songs})
}
