package user

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"../db"
	login "../login"
)

type getUser interface {
	takeUserAlbums(s db.UseDb)
	takeUserInfo(s db.UseDb)
}

type requestUser struct {
	id     string
	group  chan error
	result *User
}

// User the struct return after request /user (GET)
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

// GetAlbum the struct album of user
type GetAlbum struct {
	ID           bson.ObjectId `json:"id" bson:"_id"`
	Name         string        `json:"name" bson:"name"`
	TimeToCreate int64         `json:"timeToCreate" bson:"timeToCreate"`
	Description  string        `json:"description" bson:"description"`
}

// Info It is take user info from collections by users, images. Next step is returning User
func Info(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if id := r.Form.Get("id"); id != "" {
		albums := make([]GetAlbum, 0)
		newUser := User{Albums: albums}
		newRequestUser := requestUser{id: id, group: make(chan error), result: &newUser}
		u := getUser(newRequestUser)
		count := 0
		var s db.UseDb = &db.Session{}
		err := s.CreateSession()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer s.CloseSession()
		go u.takeUserInfo(s)
		go u.takeUserAlbums(s)
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

func (u requestUser) takeUserInfo(s db.UseDb) {
	err := s.FindUserByID(u.id, &u.result)
	u.group <- err
}

func (u requestUser) takeUserAlbums(s db.UseDb) {
	err := s.GetCollection("albums").Find(bson.M{"userId": bson.ObjectIdHex(u.id)}).All(&u.result.Albums)
	u.group <- err
}
