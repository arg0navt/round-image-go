package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"./user"
	"encoding/json"
)

type ResultIndex struct {
	Connect bool
}

func Index(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(ResultIndex{Connect: true}); err != nil {
		panic(err)
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/parse", user.ParseUsers)
	log.Fatal(http.ListenAndServe(":8080", router))
}