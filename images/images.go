package images

import (
	"encoding/base64"
	"encoding/json"
	"errors"
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

type AlbumId struct {
	ID bson.ObjectId `json:"id" bson:"_id"`
}

func CreateAlbum(w http.ResponseWriter, r *http.Request) {
	id, err := db.ValidateToken(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	var target RequestCreateAlbum
	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
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
		return
	}
	var target RequestLoadImg
	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defaultAlbumId, errDefault := foundDefaultAlbum(id)
	if errDefault != nil {
		createAlbumId, errD := createDefaultAlbum(id)
		if errD != nil {
			http.Error(w, errD.Error(), http.StatusBadRequest)
			return
		}
		defaultAlbumId = createAlbumId
	}
	json.NewEncoder(w).Encode(&defaultAlbumId)
}

func foundDefaultAlbum(id string) (AlbumId, error) {
	var getAlbum AlbumId
	err := db.GetCollection("albums").Find(bson.M{"default": true, "userId": bson.ObjectId(id)}).One(&getAlbum)
	if err != nil {
		return getAlbum, err
	}
	return getAlbum, nil
}

func createDefaultAlbum(id string) (AlbumId, error) {
	i := bson.NewObjectId()
	err := db.GetCollection("albums").Insert(bson.M{"_id": i, "userId": bson.ObjectId(id), "name": "Default album", "description": "", "timeToCreate": time.Now().Unix(), "default": true})
	if err != nil {
		return AlbumId{}, errors.New("Failed created default album")
	}
	return AlbumId{ID: i}, nil
}
