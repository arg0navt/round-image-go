package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"./user"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	defer session.Close()
	collection := session.DB("rimg").C("users")
	result := Person{}
	findUser := collection.Find(bson.M{"email": "test.test@test.com"}).One(&result)
	if findUser != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/parse", user.ParseUsers)
	log.Fatal(http.ListenAndServe(":8080", router))
}
