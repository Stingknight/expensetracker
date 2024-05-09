package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type User struct{
	Id				primitive.ObjectID	`bson:"_id"`
	Username 		*string 			`bson:"username" json:"username" validate:"required,min=2,max=100"`
	Password		*string 			`bson:"password" json:"password" validate:"required"`
	Mobile  		*string				`bson:"mobile"  json:"mobile" validate:"required,min=10"`
	Address			*string				`bson:"adderess"`
	UpdatedAt		time.Time			`bson:"updated_at"`
	CreateAt		time.Time			`bson:"created_at"`
}