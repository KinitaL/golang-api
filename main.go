package main

import (
	"myrest-api/pkg"
	"myrest-api/pkg/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	amw := middleware.AuthenticationMiddleware{}

	router.Use(amw.Middleware)

	router.HandleFunc("/register", pkg.Register).Methods("POST")
	router.HandleFunc("/login", pkg.Login).Methods("POST")
	router.HandleFunc("/logout", pkg.Logout).Methods("GET")

	router.HandleFunc("/stocks", pkg.GetNames).Methods("GET")
	router.HandleFunc("/stocks/{id}", pkg.GetName).Methods("GET")
	router.HandleFunc("/stocks", pkg.CreateName).Methods("POST")
	router.HandleFunc("/stocks/{id}", pkg.UpdateName).Methods("PUT")
	router.HandleFunc("/stocks/{id}", pkg.DeleteName).Methods("DELETE")
	http.ListenAndServe("localhost:8099", router)
}
