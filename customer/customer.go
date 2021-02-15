package customer

import (
	"context"
	"eltropy/helper"
	"eltropy/model"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

// AddCustomer - Employee can add customer to the bank
func AddCustomer(w http.ResponseWriter, r *http.Request) {

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
	if username == empID {
		// Add Customer
		var cust model.Customer
		if err := json.NewDecoder(r.Body).Decode(&cust); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}
		client, _ := helper.GetMongoClient()
		custCollection := client.Database(helper.DB).Collection(helper.CustomerCollection)
		sr := custCollection.FindOne(context.TODO(), bson.M{"custId": cust.CustId})
		if sr.Err() == nil {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(model.Response{
				Code: http.StatusConflict,
				Msg:  "customer already exist",
			})
		} else {
			_, err := custCollection.InsertOne(context.TODO(), cust)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(err)
				return
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(model.Response{
				Code: http.StatusCreated,
				Msg:  "Customer added successfully.",
			})
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "Unauthorized",
		})
	}
}

// DeleteCustomer - Employee can delete customer from bank
func DeleteCustomer(w http.ResponseWriter, r *http.Request) {

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
		// Delete customer
		var id model.Id
		json.NewDecoder(r.Body).Decode(&id)
		client, _ := helper.GetMongoClient()
		collection := client.Database(helper.DB).Collection(helper.CustomerCollection)
		sr := collection.FindOne(context.TODO(), bson.M{"custId": id.Id})
		if sr.Err() == nil {
			collection.DeleteOne(context.TODO(), bson.M{"custId": id.Id})
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(model.Response{
				Code: http.StatusOK,
				Msg:  "customer deleted. id=" + id.Id,
			})
		} else {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.Response{
				Code: http.StatusNotFound,
				Msg:  "customer not found",
			})
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "Unauthorized",
		})
	}
}
