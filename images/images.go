package images

import (
	"encoding/base64"
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

type RequestLoadImg struct {
	AlbumID string          `json:"albumId" bson:"albumId"`
	Time    int64           `json:"time" bson:"time"`
	img     base64.Encoding `json:"img" bson:"img"`
}

type Album struct {
	Name         string        `json:"name" bson:"name"`
	TimeOfCreate int64         `json:"timeOfCreate" bson:"timeOfCreate"`
	Description  string        `json:"description" bson:"description"`
	UserID       bson.ObjectId `json:"userId" bson:"userId"`
}

func CreateAlbum(w http.ResponseWriter, r *http.Request) {
	id, err := db.ValidateToken(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
	var target RequestCreateAlbum
	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
	}
	err = json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	newAlbum := Album{
		Name:         target.Name,
		TimeOfCreate: time.Now().Unix(),
		Description:  target.Description,
		UserID:       bson.ObjectId(id),
	}
	db.GetCollection("albums").Insert(&newAlbum)
	json.NewEncoder(w).Encode(&newAlbum)
}

func LoadImage(w http.ResponseWriter, r *http.Request) {
	id, err := db.ValidateToken(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
	var target RequestLoadImg
	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
	}
	err = json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	var getAlbum Album
	err = db.GetCollection("albums").Find(bson.M{"default": true, "userId": bson.ObjectIdHex(id)}).One(&getAlbum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(&getAlbum)
}
