package models

type Keyword struct{
	MovieId int `json:"movie_id" bson:"movie_id"`
	KeywordList []string `json:"keyword_list,omitempty" bson:"keyword_list,omitempty"`
}