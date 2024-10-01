package models

import (
	"time"
)

type DiscussionPart struct {
	UserId      int       `json:"user_id" bson:"user_id"`
	PartId      int       `json:"part_id" bson:"part_id"`
	Name        string    `json:"name" bson:"name"`
	ProfilePath string    `json:"profile_path" bson:"profile_path"`
	Title       string    `json:"title" bson:"title"`
	Timestamp   time.Time `json:"timestamp" bson:"timestamp"`
	Description string    `json:"description" bson:"description"`
	IsReplyOf   *int      `json:"is_reply_of,omitempty" bson:"is_reply_of, omitempty"`
}
