package main

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/sreesa7144/GoAPIExercises/controllers"
	"github.com/sreesa7144/GoAPIExercises/models"
)

func main() {

	store := sessions.NewCookieStore([]byte("secret"))
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})
	db := models.SetupModels()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})
	r.POST("/home", controllers.InsertSong)
	r.GET("/popular", controllers.PopularSongs)
	r.POST("/search", controllers.ElasticSearch)
	r.Run()
}
