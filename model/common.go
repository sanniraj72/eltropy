package model

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type User struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

type Id struct {
	Id string `json:"id"`
}
