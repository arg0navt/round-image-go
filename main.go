package main

import (
	"encoding/json"
	"log"
	"net/http"

	"./db"
	"./parse"
	user "./user"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
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
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	db.S = db.Session{session}
	defer db.S.Value.Close()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/check-token", user.CheckToken)
	router.HandleFunc("/sign-up", user.CreateUser).Methods("POST")
	router.HandleFunc("/log-in", user.LogIn).Methods("POST")
	router.HandleFunc("/parse", parse.ParseUsers)
	log.Fatal(http.ListenAndServe(":8080", router))
}
