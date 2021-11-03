package models

import "go.mongodb.org/mongo-driver/bson/primitive"

//Create Struct
type Book struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Type   string             `json:"type,omitempty" bson:"type,omitempty"`
	Name  	string             `json:"name" bson:"name,omitempty"`
	Protocol *Protocol            `json:"protocol" bson:"protocol,omitempty"`
}

type Protocol struct {
	Typee  string `json:"typee,omitempty" bson:"typee,omitempty"`
	Namee  string `json:"namee,omitempty" bson:"namee,omitempty"`
}