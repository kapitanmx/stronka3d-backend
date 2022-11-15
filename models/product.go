package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID primitive.ObjectID `bson:"_id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Imgs []string `json:"imgs"`
	Category string `json:"category"`
	Tags []string `json:"tags"`
	Price *float32 `json:"price"`
}