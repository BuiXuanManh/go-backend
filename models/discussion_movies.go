package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MovieDiscussion struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	MovieId        int                `json:"movie_id" bson:"movie_id"`
	Subject        string             `json:"subject" bson:"subject"`
	Status         bool               `json:"status" bson:"status"`
	DiscussionPart []DiscussionPart   `json:"discussion_part" bson:"discussion_part"`
}
