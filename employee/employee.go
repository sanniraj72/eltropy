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

// EmployeeSignin - Employee signin handler
func EmployeeSignin(rw http.ResponseWriter, r *http.Request) {

	rw.Header().Add("Content-Type", "application/json")
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}

	client, err := helper.GetMongoClient()
	collection := client.Database(helper.DB).Collection(helper.EmployeeCollection)
	if sr := collection.FindOne(context.TODO(), bson.M{"empId": user.Username}); sr.Err() == nil {
		// Validate username and password is correct or not
		var emp model.Employee
		if err = sr.Decode(&emp); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(rw).Encode(err)
			return
		}
		password, _ := base64.StdEncoding.DecodeString(emp.Password)
		if user.Username == emp.EmpId && user.Password == string(password) {
			// Create token if password and username is correct
			var td *helper.TokenDetails
			if td, err = helper.CreateToken(user.Username); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(rw).Encode(err)
				return
			}
			err = helper.CreateAuth(user.Username, td)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(rw).Encode(err)
				return
			}
			json.NewEncoder(rw).Encode(struct {
				Code  int    `json:"code"`
				Msg   string `json:"msg"`
				Token string `json:"token"`
			}{
				Code:  http.StatusOK,
				Msg:   "You have logged in successfully",
				Token: td.Token,
			})
		} else {
			rw.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(rw).Encode(model.Response{
				Code: http.StatusUnauthorized,
				Msg:  "username or password mismatch",
			})
		}
	} else {
		rw.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(rw).Encode(model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "username or password mismatch.",
		})
	}
}

// EmployeeSignout - Employee signout handler
func EmployeeSignout(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	helper.Signout(w, r)

}
