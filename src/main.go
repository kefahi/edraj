package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/content/{id}", NotImplementedAPI).Methods("GET")
	r.HandleFunc("/content", NotImplementedAPI).Methods("POST")
	r.HandleFunc("/content/{id}", NotImplementedAPI).Methods("PUT")
	r.HandleFunc("/content/{id}", NotImplementedAPI).Methods("DELETE")
	r.HandleFunc("/hello/", HelloAPI).Methods("GET")
	if err := http.ListenAndServe(":36794", r); err != nil {
		log.Fatal(err)
	}
}
