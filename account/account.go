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
