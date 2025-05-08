package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Moderator struct {
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name         string  `json:"name"`
    Email         string `json:"email"`
    Password      string  `json:"password"`
}
