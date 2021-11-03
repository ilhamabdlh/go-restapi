package models

import "go.mongodb.org/mongo-driver/bson/primitive"

//Create Struct
type Config struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id"`
	Type   string             `json:"type" bson:"type"`
	Name  	string             `json:"name" bson:"name"`
	Protocol *Protocol            `json:"protocol" bson:"protocol"`
}

type Protocol struct {
	Typee  string `json:"typee" bson:"typee"`
	Namee  string `json:"namee" bson:"namee"`
}