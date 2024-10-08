package models

type TopRatedMovies struct {
	Id          string  `json:"id,omitempty" bson:"_id,omitempty"`
	MovieId     int     `json:"movie_id" bson:"movie_id"`
	Title       string  `json:"title" bson:"title"`
	PosterPath  string  `json:"poster_path" bson:"poster_path"`
	VoteAvg     float64 `json:"vote_average" bson:"vote_average"`
	ReleaseDate string  `json:"release_date" bson:"release_date"`
}
