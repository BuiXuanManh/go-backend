package models

type Crew struct {
	MovieId int `json:"movie_id" bson:"movie_id"`
	Crew []string `json:"crew" bson:"crew"` 
}