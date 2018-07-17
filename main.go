package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"encoding/json"
)

type ResultIndex struct {
	Connect bool
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	result := ResultIndex{Connect:true}
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
}