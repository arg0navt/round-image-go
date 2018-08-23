package main

import (
	"encoding/json"
	"log"
	"net/http"

	"./db"
	"./user"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

type ResultIndex struct {
	Connect bool
}

type Person struct {
	ID    int    `json:"value" bson:"_id,omitempty"`
	Email string `json:"email"`
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
	var s db.SessionControll = db.Session{session}
	defer s.GetSession().Close()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/parse", user.ParseUsers)
	log.Fatal(http.ListenAndServe(":8080", router))
}
