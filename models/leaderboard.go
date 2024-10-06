package models

type Leaderboard struct {
	MovieID int `json:"movie_id" bson:"movie_id"`
	// MovieTitle     string `json:"movie_title" bson:"movie_title"`
	MoviesRated int `json:"movies_rated" bson:"movies_rated"`
	// UserId         int    `json:"user_id" bson:"user_id"`
	// Username       string `json:"username" bson:"username"`
	// PictureProfile string `json:"picture_profile" bson:"picture_profile"`
}
