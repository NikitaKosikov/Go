package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id           primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	PasswordHash string             `json:"-" bson:"password"`
	Email        string             `json:"email" bson:"email"`
	Session      Session            `json:"-" bson:"session,omitempty"`
}
