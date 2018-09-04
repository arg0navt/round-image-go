package images

import (
	"encoding/json"
	"net/http"
	"time"

	"../db"
	"gopkg.in/mgo.v2/bson"
)

type RequestCreateAlbum struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}

type Album struct {
	Name         string        `json:"name" bson:"name"`
	TimeToCreate int64         `json:"timeToCreate" bson:"timeToCreate"`
	Description  string        `json:"description" bson:"description"`
	UserID       bson.ObjectId `json:"userId" bson:"userId"`
}

func CreateAlbum(w http.ResponseWriter, r *http.Request) {
	id := db.ValidateToken(w, r)
	var target RequestCreateAlbum
	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
	}
	err := json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	newAlbum := Album{
		Name:         target.Name,
		TimeToCreate: time.Now().Unix(),
		Description:  target.Description,
		UserID:       bson.ObjectId(id),
	}
	db.GetCollection("albums").Insert(&newAlbum)
	json.NewEncoder(w).Encode(&newAlbum)
}
