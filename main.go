package main

import (
	"fmt"
	"log"
	"net/http"

	"eltropy/admin"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	fmt.Println("Starting...")
	router.HandleFunc("/admin/signup", admin.AdminSignup).Methods(http.MethodPost)
	router.HandleFunc("/admin/signin", admin.AdminSignin).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", router))
}
