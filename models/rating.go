package models

type Rating struct {
	Id        string  `json:"id,omitempty" bson:"_id,omitempty"`
	UserId    int     `json:"user_id" bson:"user_id"`
	MovieId   int     `json:"movie_id" bson:"movie_id"`
	Rating    float64 `json:"rating" bson:"rating"`
	Timestamp float64 `json:"timestamp" bson:"timestamp"`
}
