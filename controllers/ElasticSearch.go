package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	elastic "github.com/olivere/elastic/v7"
	"github.com/sreesa7144/GoAPIExercises/models"
)

func GetESClient() (*elastic.Client, error) {

	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))

	fmt.Println("ES initialized...")

	return client, err

}

func GetPopularItems(c *gin.Context) []models.Songs {

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
	return songs

}

func UniqueSearchList(duplicateSongs []models.Songs) []models.Songs {
	var songs []models.Songs

	uniqueMap := make(map[uint]bool)

	for _, song := range duplicateSongs {
		if _, value := uniqueMap[song.SongID]; !value {
			uniqueMap[song.SongID] = true
			songs = append(songs, song)
		}
	}
	return songs
}

func ElasticSearch(c *gin.Context) {
	session := sessions.Default(c)

	ctx := context.Background()

	esclient, _ := GetESClient()
	resp := make(map[string]string)
	var songs []models.Songs
	var popularSongs []models.Songs

	if err := c.ShouldBindJSON(&resp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMultiMatchQuery("*"+resp["name"]+"*", "title", "movie", "language").Type("phrase_prefix"))
	searchService := esclient.Search().Index("songs").SearchSource(searchSource)
	searchResult, err := searchService.Do(ctx)
	if err != nil {
		fmt.Println("[ProductsES][GetPIds]Error=", err)
		return
	}
	fmt.Println("searchResults  len is", len(searchResult.Hits.Hits))

	for _, hit := range searchResult.Hits.Hits {
		var song models.Songs
		err := json.Unmarshal(hit.Source, &song)
		fmt.Println("Unmarshalled song is", song)
		if err != nil {
			fmt.Println("[Getting Songs][Unmarshal] Err=", err)
		}

		songs = append(songs, song)
	}
	if songs == nil {
		songs = GetPopularItems(c)
		log.Println("Popular Songs is  ", songs)
	} else {

		log.Println("Songs is  ", songs)
		v := session.Get("searched")
		if v == nil {
			popularSongs = GetPopularItems(c)
			log.Println("Popular Songs is (func)  ", popularSongs)
			for i := 0; i < len(popularSongs); i += 1 {
				if i < len(songs) {

					popularSongs[i] = songs[i]

				} else {
					log.Println("Popular Songs after is", popularSongs)
					break
				}
			}

		} else {

			str := v.(string)
			_ = json.Unmarshal([]byte(str), &popularSongs)
			log.Println("Cookie Songs is   ", popularSongs)
			for i := 0; i < len(popularSongs); i += 1 {
				if i < len(songs) {
					popularSongs[i] = songs[i]
				} else {
					break
				}
			}
		}

		songs = popularSongs
	}
	jsonDat, _ := json.Marshal(songs)
	session.Set("searched", string(jsonDat))
	session.Save()
	c.JSON(http.StatusOK, gin.H{"data": UniqueSearchList(songs)})

}
