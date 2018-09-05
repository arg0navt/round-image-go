package getUser

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"../db"
	user "../user"
)

// type UserInterface interface {
// 	takeUserInfo(c chan int, id string)
// }

type GetUser struct {
	ID              bson.ObjectId `json:"id" bson:"_id"`
	FirstName       string        `json:"first_name" bson:"first_name"`
	LastName        string        `json:"last_name" bson:"last_name"`
	Email           string        `json:"email" bson:"email"`
	Verification    bool          `json:"verification" bson:"verification"`
	DateLastActive  int64         `json:"dateLastActive" bson:"dateLastActive"`
	user.DetailInfo `json:"detailInfo" bson:"detailInfo"`
	Albums          []GetAlbum `json:"albums" bson:"albums"`
}

type GetAlbum struct {
	ID           bson.ObjectId `json:"id" bson:"_id"`
	Name         string        `json:"name" bson:"name"`
	TimeToCreate int64         `json:"timeToCreate" bson:"timeToCreate"`
	Description  string        `json:"description" bson:"description"`
}

func UserInfo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var u GetUser
	if id := r.Form.Get("id"); id != "" {
		count := 0
		group := make(chan string)
		go u.takeUserInfo(group, id)
		go u.takeUserAlbums(group, id)
		for count <= 1 {
			select {
			case result := <-group:
				if result != "" {
					http.Error(w, result, http.StatusInternalServerError)
					return
				}
				count++
			}
		}
		json.NewEncoder(w).Encode(u)
	}
}

func (u *GetUser) takeUserInfo(c chan string, id string) {
	err := db.GetCollection("users").FindId(bson.ObjectIdHex(id)).One(&u)
	if err != nil {
		c <- "user not found"
	}
	c <- ""
}

func (u *GetUser) takeUserAlbums(c chan string, id string) {
	err := db.GetCollection("albums").Find(bson.M{"userId": bson.ObjectIdHex(id)}).All(&u.Albums)
	if err != nil {
		c <- "error get user albums"
	}
	c <- ""
}
