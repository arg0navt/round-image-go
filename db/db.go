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

var S Session

func GetUsers() *mgo.Collection {
	return S.Value.DB(DB).C("users")
}

func FindUser(email string) bson.ObjectId {
	var result bson.ObjectIdвыа
	S.Value.DB(DB).C("users").Find(bson.M{"email": email})
	return result
}
