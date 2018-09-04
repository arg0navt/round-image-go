package getUser

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	db "../db"
	user "../user"
)

type GetUser struct {
	FirstName       string       `json:"first_name" bson:"first_name"`
	LastName        string       `json:"last_name" bson:"last_name"`
	Email           string       `json:"email" bson:"email"`
	Verification    bool         `json:"verification" bson:"verification"`
	DateLastActive  int64        `json:"dateLastActive" bson:"dateLastActive"`
	Albums          []user.Album `json:"albums" bson:"albums"`
	user.DetailInfo `json:"detailInfo" bson:"detailInfo"`
}

func UserInfo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if id := r.Form.Get("id"); id != "" {
		var result GetUser
		err := db.CollectionUsers().Find(bson.M{"id": bson.ObjectId(id)}).One(&result)
		fmt.Println(err)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(result)
	}
}
