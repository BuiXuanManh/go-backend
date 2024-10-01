package models

type Genre struct {
	ID   int    `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}