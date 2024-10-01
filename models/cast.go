package models

type Cast struct{
	MovieId int `json:"movie_id" bson:"movie_id"`
	Cast []string  `json:"cast" bson:"cast"`
}