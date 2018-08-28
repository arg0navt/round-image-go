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

func ThereIsUser(email string) bool {
	result, _ := S.Value.DB(DB).C("users").Find(bson.M{"email": email}).Count()
	if result != 0 {
		return true
	}
	return false
}
