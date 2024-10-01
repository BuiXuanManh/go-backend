package models

type User struct {
	UserId             int    `json:"id" bson:"id"`
	Username           string `json:"username" bson:"username"`
	PasswordHash       string `json:"password_hash" bson:"password_hash"`
	FavoriteList       []int  `json:"favorite_list" bson:"favorite_list"`
	RecommendationList []int  `json:"recommendation_list" bson:"recommendation_list"`
	WatchList          []int  `json:"watch_list" bson:"watch_list"`
	// Settings           []Document `json:"settings"`
	PictureProfile string `json:"picture_profile" bson:"picture_profile"`
}
