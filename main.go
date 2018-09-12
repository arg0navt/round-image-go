package main

import (
	"encoding/json"
	"log"
	"net/http"

	"./images"
	"./login"
	"./user"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

// ResultIndex main page struct
type ResultIndex struct {
	Connect bool
}

// Index return json connect
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
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/check_token", login.CheckToken)
	router.HandleFunc("/sign_up", login.CreateUser).Methods("POST")
	router.HandleFunc("/log_out", login.LogOut)
	router.HandleFunc("/log_in", login.LogIn).Methods("POST")
	router.HandleFunc("/user", user.Info)
	router.HandleFunc("/create_album", images.CreateAlbum).Methods("POST")
	router.HandleFunc("/load_image", images.LoadImage).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
