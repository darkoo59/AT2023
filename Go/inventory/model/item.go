package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Item struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string
	Quantity uint32
	Price    float64
}
