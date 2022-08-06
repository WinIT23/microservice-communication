package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Comment struct {
	Id     primitive.ObjectID `json:"id"`
	PostId primitive.ObjectID `json:"postid"`
	Text   string             `json:"text"`
}
