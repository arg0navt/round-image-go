package main

import (
	"encoding/json"
	"log"
	"net/http"

	"./db"
	"./getUser"
	"./images"
	"./parse"
	"./user"
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
	router.HandleFunc("/check_token", user.CheckToken)
	router.HandleFunc("/sign_up", user.CreateUser).Methods("POST")
	router.HandleFunc("/log_out", user.LogOut)
	router.HandleFunc("/log_in", user.LogIn).Methods("POST")
	router.HandleFunc("/parse", parse.ParseUsers)
	router.HandleFunc("/user", getUser.UserInfo)
	router.HandleFunc("/create_album", images.CreateAlbum).Methods("POST")
	router.HandleFunc("/load_image", images.LoadImage).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
