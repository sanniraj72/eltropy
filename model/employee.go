package model

type Employee struct {
	EmpId    string `bson:"empId"`
	Password string `bson:"password"`
	Name     string `bson:"name"`
	Address  string `bson:"address"`
	Phone    string `bson:"phone"`
	Email    string `bson:"email"`
}
