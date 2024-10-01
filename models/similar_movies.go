package models

type SimilarMovies struct{
	MovieId int `json:"id" bson:"id"`
	SimilarMovies []int `json:"similar_movies" bson:"similar_movies" `
}

