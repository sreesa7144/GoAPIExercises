package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sreesa7144/GoAPIExercises/models"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
}

func InsertSong(c *gin.Context) {

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-south-1"),
		Credentials: credentials.NewStaticCredentials("", "", ""),
	})

	uploader := s3manager.NewUploader(sess)

	db := c.MustGet("db").(*gorm.DB)
	ctx := context.Background()
	esclient, err := GetESClient()
	if err != nil {
		fmt.Println("Error initializing : ", err)

	}

	title := c.PostForm("title")
	movie := c.PostForm("movie")
	language := c.PostForm("language")
	icon := c.PostForm("icon")
	rating := c.PostForm("popularity")

	file, _ := c.FormFile("song")
	log.Println(file.Filename)
	out, _ := file.Open()
	content, _ := ioutil.ReadAll(out)
	err = ioutil.WriteFile("C:/Users/panchangam/go/src/github.com/sreesa7144/GoAPIExercises/files/"+file.Filename, content, 0644)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("songsstorage"),
		Key:    aws.String(file.Filename),
		Body:   out,
	})
	if err != nil {
		// Print the error and exit.
		exitErrorf("Unable to upload %q to songsstorage, %v", file.Filename, err)
	} else {

		log.Println(result)
		fmt.Printf("Successfully uploaded %q to songsstorage\n", file.Filename)
	}

	if err != nil {
		log.Println(err)
	}

	popularity, _ := strconv.ParseInt(rating, 10, 64)
	insertSong := models.Songs{Title: title, Movie: movie, Language: language, Icon: icon, SongURI: result.Location, Popularity: popularity}
	db.Create(&insertSong)
	dataJSON, err := json.Marshal(insertSong)
	js := string(dataJSON)
	ind, err := esclient.Index().
		Index("songs").
		BodyJson(js).
		Do(ctx)

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": insertSong, "ElasticData": ind})
}
