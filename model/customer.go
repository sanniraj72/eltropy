package model

type Customer struct {
	Name       string    `bson:"name"`
	CustId     string    `bson:"custId"`
	Address    string    `bson:"address"`
	Phone      string    `bson:"phone"`
	Email      string    `bson:"email"`
	BranchName string    `bson:"branchName"`
	BranchAdd  string    `bson:"branchAdd"`
	KYC        bool      `bson:"kyc"`
	Accounts   []Account `bson:"accounts"`
}
