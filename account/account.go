package account

import (
	"context"
	"eltropy/helper"
	"eltropy/model"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

// CreateAccount - Employee can create account and link with customer
func CreateAccount(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	// Extract token
	ad, err := helper.ExtractToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "Unauthorized",
		})
		return
	}

	// Fetch auth from redis
	username, err := helper.FetchAuth(ad)
	if username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "Unauthorized",
		})
		return
	}

	values := r.URL.Query()
	empID := values.Get("empId")
	custID := values.Get("custId")
	if empID == username {
		// Create account
		client, _ := helper.GetMongoClient()
		collection := client.Database(helper.DB).Collection(helper.CustomerCollection)
		sr := collection.FindOne(context.TODO(), bson.M{"custId": custID})
		if sr.Err() == nil {
			// Link account with customer
			var customer model.Customer
			var account model.Account
			sr.Decode(&customer)
			err = json.NewDecoder(r.Body).Decode(&account)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(err)
				return
			}
			// Check for duplicate account
			acctFound := false
			if customer.Accounts != nil {
				for _, acct := range customer.Accounts {
					if acct.Number == account.Number {
						acctFound = true
						break
					}
				}
			}
			if acctFound {
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(model.Response{
					Code: http.StatusConflict,
					Msg:  "account already exist",
				})
				return
			}
			customer.Accounts = append(customer.Accounts, account)
			result, err := collection.UpdateOne(
				context.TODO(),
				bson.M{"custId": custID},
				bson.D{
					{"$set", bson.D{{"accounts", customer.Accounts}}},
				},
			)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(model.Response{
					Code: http.StatusInternalServerError,
					Msg:  "Account couldn't link with customer",
				})
				return
			}
			if result.ModifiedCount == 1 {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(model.Response{
					Code: http.StatusOK,
					Msg:  "Account created and linked successfully.",
				})
				return
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.Response{
				Code: http.StatusNotFound,
				Msg:  "Customer not found",
			})
			return
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "empId required",
		})
		return
	}
}

// GetAccountBalance - get account balance of an account
func GetAccountBalance(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	// Extract token
	ad, err := helper.ExtractToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "Unauthorized",
		})
		return
	}

	// Fetch auth from redis
	username, err := helper.FetchAuth(ad)
	if username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "Unauthorized",
		})
		return
	}

	values := r.URL.Query()
	empID := values.Get("empId")

	if empID == username {
		custID := values.Get("custId")
		accountNumber := values.Get("acctId")
		// get customer
		client, _ := helper.GetMongoClient()
		collection := client.Database(helper.DB).Collection(helper.CustomerCollection)
		sr := collection.FindOne(context.TODO(), bson.M{"custId": custID})
		if sr.Err() != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.Response{
				Code: http.StatusNotFound,
				Msg:  "customer not found",
			})
			return
		}
		var customer model.Customer
		sr.Decode(&customer)
		accounts := customer.Accounts
		if accounts == nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.Response{
				Code: http.StatusNotFound,
				Msg:  "Customer has no account.",
			})
			return
		}
		for _, account := range accounts {
			if account.Number == accountNumber {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(account)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusNotFound,
			Msg:  "Account not found",
		})
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "unauthorized",
		})
	}
}

// TransferMoney - Transfer money
func TransferMoney(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	// Extract token
	ad, err := helper.ExtractToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "Unauthorized",
		})
		return
	}

	// Fetch auth from redis
	username, err := helper.FetchAuth(ad)
	if username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "Unauthorized",
		})
		return
	}

	values := r.URL.Query()
	empID := values.Get("empId")
	if empID == username {
		// Transfer process
		var transfer model.Transfer
		json.NewDecoder(r.Body).Decode(&transfer)
		client, _ := helper.GetMongoClient()
		collection := client.Database(helper.DB).Collection(helper.CustomerCollection)
		srcSR := collection.FindOne(context.TODO(), bson.M{"custId": transfer.SrcCustomer})
		destSR := collection.FindOne(context.TODO(), bson.M{"custId": transfer.DestCustomer})
		if srcSR.Err() != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.Response{
				Code: http.StatusNotFound,
				Msg:  "customer not found, id-" + transfer.SrcCustomer,
			})
			return
		}
		if destSR.Err() != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.Response{
				Code: http.StatusNotFound,
				Msg:  "customer not found, id-" + transfer.DestCustomer,
			})
			return
		}
		var srcCustomer, destCustomer model.Customer
		srcSR.Decode(&srcCustomer)
		destSR.Decode(&destCustomer)
		if srcCustomer.Accounts == nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.Response{
				Code: http.StatusNotFound,
				Msg:  "customer do not have any account, custId-" + transfer.SrcCustomer,
			})
			return
		}
		if destCustomer.Accounts == nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.Response{
				Code: http.StatusNotFound,
				Msg:  "customer do not have any account, custId-" + transfer.DestCustomer,
			})
			return
		}
		// Deduct amount from source
		for i, acct := range srcCustomer.Accounts {
			if acct.Number == transfer.SrcAccount {
				srcCustomer.Accounts[i].Balance = acct.Balance - transfer.Amount
				collection.UpdateOne(
					context.TODO(),
					bson.M{"custId": transfer.SrcCustomer},
					bson.D{
						{"$set", bson.D{{"accounts", srcCustomer.Accounts}}},
					},
				)
				break
			}
		}
		// Add amount to destination
		for i, acct := range destCustomer.Accounts {
			if acct.Number == transfer.DestAccount {
				destCustomer.Accounts[i].Balance = acct.Balance + transfer.Amount
				collection.UpdateOne(
					context.TODO(),
					bson.M{"custId": transfer.DestCustomer},
					bson.D{
						{"$set", bson.D{{"accounts", destCustomer.Accounts}}},
					},
				)
				break
			}
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusOK,
			Msg:  "Transferred suucessfully.",
		})
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "unauthorized",
		})
	}
}
