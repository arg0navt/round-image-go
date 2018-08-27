package db

import (
	"gopkg.in/mgo.v2"
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
