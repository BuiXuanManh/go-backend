package models

type UserReview struct {
	UserID         int    `bson:"user_id" json:"user_id"`
	MovieID        int    `bson:"movie_id" json:"movie_id"`
	Username       string `bson:"username" json:"username"`
	PictureProfile string `bson:"picture_profile" json:"picture_profile"`
	Comment        string `bson:"comment" json:"comment"`
}
