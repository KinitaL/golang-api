package main

import (
	"myrest-api/pkg"
	"myrest-api/pkg/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	amw := middleware.AuthenticationMiddleware{TokenUsers: make(map[string]string)}
	amw.Populate()
	router.Use(amw.Middleware)

	router.HandleFunc("/people", pkg.GetNames).Methods("GET")
	router.HandleFunc("/people/{id}", pkg.GetName).Methods("GET")
	router.HandleFunc("/people", pkg.CreateName).Methods("POST")
	router.HandleFunc("/people/{id}", pkg.UpdateName).Methods("PUT")
	router.HandleFunc("/people/{id}", pkg.DeleteName).Methods("DELETE")
	http.ListenAndServe("localhost:8099", router)
}
