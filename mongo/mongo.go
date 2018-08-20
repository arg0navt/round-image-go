package mongo

import (
	"gopkg.in/mgo.v2"
)

type ConnectDb struct {
	connect *mgo.Session
}

func NewSession(url string) (*ConnectDb,error) {
	connect, err := mgo.Dial(url)
	if err != nil {
	  return nil,err
	}
	return &ConnectDb{connect}, err
}

func(s *ConnectDb) GetCollection(db string, col string) *mgo.Collection {
	return s.connect.DB(db).C(col)
}

func(c *ConnectDb) Close() {
	if(c.connect != nil) {
	  c.connect.Close()
	}
}