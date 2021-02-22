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
}

type Transfer struct {
	SrcAccount   string  `json:"srcAccount"`
	SrcCustomer  string  `json:"srcCustomer"`
	DestAccount  string  `json:"destAccount"`
	DestCustomer string  `json:"destCustomer"`
	Amount       float64 `json:"amount"`
}
