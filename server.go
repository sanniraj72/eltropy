package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/eltropy/admin"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt/options"
)

func main() {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)

	router := mux.NewRouter()
	router.HandleFunc("/admin/signup", admin.AdminSignup)
	http.ListenAndServe(":8080", router)
}
