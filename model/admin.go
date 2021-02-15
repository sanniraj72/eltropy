package model

type Admin struct {
	Name     string `bson:"name"`
	Address  string `bson:"address"`
	Phone    string `bson:"phone"`
	Email    string `bson:"email"`
	UserName string `bson:"username"`
	Password string `bson:"password"`
}
