package admin

import (
	"context"
	"eltropy/helper"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"eltropy/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// AdminSignup - Admin signup handler
func AdminSignup(rw http.ResponseWriter, r *http.Request) {

	rw.Header().Set("Content-Type", "application/json")
	var admin model.Admin
	err := json.NewDecoder(r.Body).Decode(&admin)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}

	// check for duplicate username
	client, err := helper.GetMongoClient()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode(err)
		return
	}
	collection := client.Database(helper.DB).Collection(helper.AdminCollection)
	sr := collection.FindOne(context.TODO(), bson.M{"username": admin.UserName})
	if sr.Err() == mongo.ErrNoDocuments {
		// Create new entry
		admin.Password = base64.StdEncoding.EncodeToString([]byte(admin.Password))
		_, err = collection.InsertOne(context.TODO(), admin)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(rw).Encode(err)
			return
		}
		rw.WriteHeader(http.StatusCreated)
		json.NewEncoder(rw).Encode(model.Response{
			Code: http.StatusCreated,
			Msg:  "Registered as an admin Successfully.",
		})
	} else {
		rw.WriteHeader(http.StatusConflict)
		json.NewEncoder(rw).Encode(model.Response{
			Code: http.StatusConflict,
			Msg:  "Admin user already exist",
		})
	}
}

// AdminSignin - Admin signin handler
func AdminSignin(rw http.ResponseWriter, r *http.Request) {

	rw.Header().Add("Content-Type", "application/json")
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}

	client, err := helper.GetMongoClient()
	collection := client.Database(helper.DB).Collection(helper.AdminCollection)
	if sr := collection.FindOne(context.TODO(), bson.M{"username": user.Username}); sr.Err() == nil {
		// Validate username and password is correct or not
		var admin model.Admin
		if err = sr.Decode(&admin); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(rw).Encode(err)
			return
		}
		password, _ := base64.StdEncoding.DecodeString(admin.Password)
		if user.Username == admin.UserName && user.Password == string(password) {
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
		rw.WriteHeader(http.StatusForbidden)
		json.NewEncoder(rw).Encode(model.Response{
			Code: http.StatusForbidden,
			Msg:  "You dont't have admin access.",
		})
	}
}

// AdminSignout - Admin signout handler
func AdminSignout(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	helper.Signout(w, r)
}
