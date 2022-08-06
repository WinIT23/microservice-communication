package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Comment struct {
	Id     primitive.ObjectID `json:"id"`
	PostId primitive.ObjectID `json:"post_id"`
	Text   string             `json:"text"`
}
