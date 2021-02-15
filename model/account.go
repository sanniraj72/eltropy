package model

type Account struct {
	Number       string        `bson:"number"`
	IFSC         string        `bson:"ifsc"`
	Type         string        `bson:"type"`
	Balance      float64       `bson:"balance"`
	OpeningDate  string        `bson:"openingDate"`
	ClosingDate  string        `bson:"closingDate"`
	Status       string        `bson:"status"` // Active/Inactive/Blocked
	Transactions []Transaction `bson:"transactions"`
}

type Transaction struct {
	Date   string  `bson:"date"`
	Amount float64 `bson:"amount"`
	Type   string  `bson:"type"` // debited/credited
}
