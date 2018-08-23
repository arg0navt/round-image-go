package db

import (
	"gopkg.in/mgo.v2"
)

const URL = "localhost:27017"
const DB = "rimg"

type SessionControll interface {
	GetSession()
	GetUsers()
}

type Session struct {
	Value *mgo.Session
}

func (s Session) GetSession() *mgo.Session {
	return s.Value
}

func (s Session) GetUsers() *mgo.Collection {
	collection := s.Value.DB("rimg").C("users")
	return collection
}
