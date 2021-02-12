package admin

import (
	"context"
	"eltropy/db"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"eltropy/model"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
)

type response struct {
	code int
	msg  string
}

func AdminSignup(rw http.ResponseWriter, r *http.Request) {

	rw.Header().Add("Content-Type", "application/json")
	var admin model.Admin
	err := json.NewDecoder(r.Body).Decode(&admin)
	if err != nil {
		log.Fatal(err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	// check for duplicate username
	client, err := db.GetMongoClient()
	if err != nil {
		log.Fatal(err)
		json.NewEncoder(rw).Encode(err)
		return
	}
	collection := client.Database(db.DB).Collection(db.ADMIN_COLLECTION)
	var existingAdmin model.Admin
	err = collection.FindOne(context.TODO(), bson.M{"username": admin.UserName}).Decode(&existingAdmin)
	if err != nil {
		log.Fatal(err)
		json.NewEncoder(rw).Encode(err)
		return
	}
	if admin.UserName == existingAdmin.UserName {
		json.NewEncoder(rw).Encode(response{
			code: http.StatusConflict,
			msg:  "Admin user already exist",
		})
		return
	}
	admin.Password = base64.StdEncoding.EncodeToString([]byte(admin.Password))
	if err != nil {
		json.NewEncoder(rw).Encode(err)
		return
	}
	_, err = collection.InsertOne(context.TODO(), admin)
	if err != nil {
		json.NewEncoder(rw).Encode(err)
		return
	}
	p := struct {
		code int
		msg  string
	}{
		code: http.StatusCreated,
		msg:  "Registered as an admin Successfully.",
	}
	json.NewEncoder(rw).Encode(p)
}

func AdminSignin(rw http.ResponseWriter, r *http.Request) {

	rw.Header().Add("Content-Type", "application/json")
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
		json.NewEncoder(rw).Encode(err)
		return
	}

	// Validate username and password is correct or not
	client, err := db.GetMongoClient()
	collection := client.Database(db.DB).Collection(db.ADMIN_COLLECTION)
	var admin model.Admin
	if err = collection.FindOne(context.TODO(), bson.M{}).Decode(&admin); err != nil {
		log.Fatal(err)
	}
	b, err := base64.StdEncoding.DecodeString(admin.Password)
	if err != nil {
		log.Fatal(err)
		json.NewEncoder(rw).Encode(err)
		return
	}
	if user.Username != admin.UserName || user.Password != string(b) {
		json.NewEncoder(rw).Encode(struct {
			msg string
		}{
			msg: "Wrong Username or password.",
		})
		return
	}

	// Create token if password and username is correct
	token, err := createToken(user.Username)
	json.NewEncoder(rw).Encode(struct {
		code  int
		token string
	}{
		code:  http.StatusOK,
		token: token,
	})
}

func createToken(username string) (string, error) {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = username
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return token, nil
}
