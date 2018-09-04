package getUser

import (
	"fmt"
	"net/http"

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
	for key, value := range r.Form {
		fmt.Println(key, ": ", value)
	}
}
