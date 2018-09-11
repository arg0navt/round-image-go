package user

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"../db"
	login "../login"
)

type getUser interface {
	takeUserAlbums()
	takeUserInfo()
}

type requestUser struct {
	id     string
	group  chan error
	result *User
}

type User struct {
	ID               bson.ObjectId `json:"id" bson:"_id"`
	FirstName        string        `json:"first_name" bson:"first_name"`
	LastName         string        `json:"last_name" bson:"last_name"`
	Email            string        `json:"email" bson:"email"`
	Verification     bool          `json:"verification" bson:"verification"`
	DateLastActive   int64         `json:"dateLastActive" bson:"dateLastActive"`
	login.DetailInfo `json:"detailInfo" bson:"detailInfo"`
	Albums           []GetAlbum `json:"albums" bson:"albums"`
}

type GetAlbum struct {
	ID           bson.ObjectId `json:"id" bson:"_id"`
	Name         string        `json:"name" bson:"name"`
	TimeToCreate int64         `json:"timeToCreate" bson:"timeToCreate"`
	Description  string        `json:"description" bson:"description"`
}

func UserInfo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if id := r.Form.Get("id"); id != "" {
		albums := make([]GetAlbum, 0)
		newUser := User{Albums: albums}
		newRequestUser := requestUser{id: id, group: make(chan error), result: &newUser}
		u := getUser(newRequestUser)
		count := 0
		go u.takeUserInfo()
		go u.takeUserAlbums()
		for count <= 1 {
			select {
			case err := <-newRequestUser.group:
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				count++
			}
		}
		json.NewEncoder(w).Encode(newRequestUser.result)
	}
}

func (u requestUser) takeUserInfo() {
	err := db.GetCollection("users").FindId(bson.ObjectIdHex(u.id)).One(&u.result)
	u.group <- err
}

func (u requestUser) takeUserAlbums() {
	err := db.GetCollection("albums").Find(bson.M{"userId": bson.ObjectIdHex(u.id)}).All(&u.result.Albums)
	u.group <- err
}
