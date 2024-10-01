package models

type SpokenLanguage struct {
	EnglishName string `json:"english_name" bson:"english_name"`
	Iso639_1    string `json:"iso_639_1" bson:"iso_639_1"`
	Name        string `json:"name" bson:"name"`
}