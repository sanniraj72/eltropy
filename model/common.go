package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Config struct {
	Uri string
}

type User struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

type Admin struct {
	Name     string `bson:name`
	Address  string `bson:address`
	Phone    string `bson:phone`
	Email    string `bson:email`
	UserName string `bson:"username"`
	Password string `bson:"password"`
}

type AdminResponse struct {
	_ID      primitive.ObjectID `bson:"_id"`
	Token    string             `bson:""token`
	Username string
}
