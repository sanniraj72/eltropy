package main

import (
	"fmt"
	"log"
	"net/http"

	"eltropy/account"
	"eltropy/admin"
	"eltropy/customer"
	"eltropy/employee"
	"eltropy/helper"

	"github.com/gorilla/mux"
)

func main() {

	helper.InistalizeRedis()
	router := mux.NewRouter()
	fmt.Println("Starting...")
	router.HandleFunc("/admin/signup", admin.AdminSignup).Methods(http.MethodPost)
	router.HandleFunc("/admin/signin", admin.AdminSignin).Methods(http.MethodPost)
	router.HandleFunc("/admin/signout", admin.AdminSignout).Methods(http.MethodPost)
	router.HandleFunc("/employee/add", employee.AddEmployee).Methods(http.MethodPost)
	router.HandleFunc("/employee/delete", employee.DeleteEmployee).Methods(http.MethodDelete)
	router.HandleFunc("/employee/signin", employee.EmployeeSignin).Methods(http.MethodPost)
	router.HandleFunc("/employee/signout", employee.EmployeeSignout).Methods(http.MethodPost)
	router.HandleFunc("/customer/add", customer.AddCustomer).Methods(http.MethodPost)
	router.HandleFunc("/customer/delete", customer.DeleteCustomer).Methods(http.MethodDelete)
	router.HandleFunc("/account/add", account.CreateAccount).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", router))
}
