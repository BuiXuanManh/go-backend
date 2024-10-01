package models

type Leaderboard struct{
	UserId int `json:"user_id" bson:"user_id"`
	MoviesRated int `json:"movies_rated" bson:"movies_rated"`
	Username string `json:"username" bson:"username"`
	PictureProfile string `json:"picture_profile" bson:"picture_profile"`	
}

