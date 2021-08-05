package main

import (
	"myrest-api/pkg"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/people", pkg.GetNames).Methods("GET")
	router.HandleFunc("/people/{id}", pkg.GetName).Methods("GET")
	router.HandleFunc("/people", pkg.CreateName).Methods("POST")
	http.ListenAndServe("localhost:8099", router)
}
