package employee

import (
	"context"
	"eltropy/helper"
	"eltropy/model"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

// AddEmployee - Admin can add employee to the bank
func AddEmployee(w http.ResponseWriter, r *http.Request) {

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

	// Add Employee
	var emp model.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	client, _ := helper.GetMongoClient()
	collection := client.Database(helper.DB).Collection(helper.EmployeeCollection)
	sr := collection.FindOne(context.TODO(), bson.M{"empId": emp.EmpId})
	if sr.Err() == nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusConflict,
			Msg:  "employee already exist",
		})
	} else {
		emp.Password = base64.StdEncoding.EncodeToString([]byte(emp.Password))
		_, err := collection.InsertOne(context.TODO(), emp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusCreated,
			Msg:  "Employee added successfully.",
		})
	}
}

// DeleteEmployee - Admin can delete employee to the bank
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {

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

	// Delete employee
	var id model.Id
	json.NewDecoder(r.Body).Decode(&id)
	client, _ := helper.GetMongoClient()
	collection := client.Database(helper.DB).Collection(helper.EmployeeCollection)
	sr := collection.FindOne(context.TODO(), bson.M{"empId": id.Id})
	if sr.Err() == nil {
		collection.DeleteOne(context.TODO(), bson.M{"empId": id.Id})
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusOK,
			Msg:  "employee deleted. id=" + id.Id,
		})
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusNotFound,
			Msg:  "employee not found",
		})
	}
}
