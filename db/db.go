package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const URL = "localhost:27017"
const DB = "rimg"

type Session struct {
	Value *mgo.Session
}

type UserID struct {
	ID bson.ObjectId `json:"id" bson:"_id"`
}

var S Session

func GetCollection(name string) *mgo.Collection {
	return S.Value.DB(DB).C(name)
}

func ThereIsUserEmail(email string) bool {
	result, _ := GetCollection("users").Find(bson.M{"email": email}).Count()
	if result != 0 {
		return true
	}
	return false
}

func GetUserId(email string) string {
	var result UserID
	err := GetCollection("users").Find(bson.M{"email": email}).One(&result)
	if err != nil {
		return ""
	}
	return string(result.ID)
}
