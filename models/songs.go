package models

type Songs struct {
	SongID     uint   `json:"id" gorm:"primary_key"`
	Title      string `json:"title"`
	Movie      string `json:"movie"`
	Language   string `json:"language"`
	Icon       string `json:"icon"`
	SongURI    string `json:"songuri"`
	Popularity int64  `json:"popularity"`
}

type SongSearch struct {
	Title    string `json:"title"`
	Movie    string `json:"movie"`
	Language string `json:"language"`
}
