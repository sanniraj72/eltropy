package admin

import (
	"context"
	"eltropy/db"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"eltropy/model"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

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
	client, err := db.GetMongoClient()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode(err)
		return
	}
	collection := client.Database(db.DB).Collection(db.ADMIN_COLLECTION)
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
		json.NewEncoder(rw).Encode(response{
			Code: http.StatusCreated,
			Msg:  "Registered as an admin Successfully.",
		})
	} else {
		rw.WriteHeader(http.StatusConflict)
		json.NewEncoder(rw).Encode(response{
			Code: http.StatusConflict,
			Msg:  "Admin user already exist",
		})
	}
}

func AdminSignin(rw http.ResponseWriter, r *http.Request) {

	rw.Header().Add("Content-Type", "application/json")
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	client, err := db.GetMongoClient()
	collection := client.Database(db.DB).Collection(db.ADMIN_COLLECTION)
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
			var token string
			if token, err = createToken(user.Username); err != nil {
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
				Token: token,
			})
		} else {
			rw.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(rw).Encode(response{
				Code: http.StatusUnauthorized,
				Msg:  "username or password mismatch",
			})
		}
	} else {
		rw.WriteHeader(http.StatusForbidden)
		json.NewEncoder(rw).Encode(response{
			Code: http.StatusForbidden,
			Msg:  "You dont't have admin access.",
		})
	}
}

func createToken(username string) (string, error) {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["username"] = username
	atClaims["expiry"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return token, nil
}
