package images

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"../db"
)

type Album struct {
	Name         string        `json:"name" bson:"name"`
	TimeToCreate int64         `json:"timeToCreate" bson:"timeToCreate"`
	Description  string        `json:"description" bson:"description"`
	UserID       bson.ObjectId `json:"userId" bson:"userId"`
}

func createAlbum(name string, id bson.ObjectId) {
	newAlbum := Album{
		Name:         name,
		TimeToCreate: time.Now().Unix(),
		UserID:       id,
	}
	db.GetCollection("albums").Insert(&newAlbum)
}
