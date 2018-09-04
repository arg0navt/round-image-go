package getUser

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	db "../db"
	user "../user"
)

type GetUser struct {
	FirstName       string `json:"first_name" bson:"first_name"`
	LastName        string `json:"last_name" bson:"last_name"`
	Email           string `json:"email" bson:"email"`
	Verification    bool   `json:"verification" bson:"verification"`
	DateLastActive  int64  `json:"dateLastActive" bson:"dateLastActive"`
	user.DetailInfo `json:"detailInfo" bson:"detailInfo"`
}

type UserID struct {
	ID bson.ObjectId `json:"id" bson:"_id"`
}

func UserInfo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if id := r.Form.Get("id"); id != "" {
		var result GetUser
		err := db.GetCollection("users").FindId(bson.ObjectIdHex(id)).One(&result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(result)
	}
}

func GetUserId(email string) string {
	var result UserID
	err := db.GetCollection("users").Find(bson.M{"email": email}).One(&result)
	if err != nil {
		return ""
	}
	return string(result.ID)
}
