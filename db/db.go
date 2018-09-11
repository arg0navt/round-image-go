package db

import (
	"errors"
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const URL = "localhost:27017"
const DB = "rimg"

type UserID struct {
	ID bson.ObjectId `json:"id" bson:"_id"`
}

type UseDb interface {
	CreateSession() error
	CloseSession()
	GetCollection(name string) *mgo.Collection
}

type Session struct {
	Value *mgo.Session
}

func (s Session) CreateSession() error {
	connect, err := mgo.Dial(URL)
	if err != nil {
		return err
	}
	s.Value = connect
	return nil
}

func (s Session) CloseSession() {
	s.Value.Close()
}

func (s Session) GetCollection(name string) *mgo.Collection {
	return s.Value.DB(DB).C(name)
}

func ThereIsUserEmail(s UseDb, email string) bool {
	result, _ := s.GetCollection("users").Find(bson.M{"email": email}).Count()
	if result != 0 {
		return true
	}
	return false
}

func GetUserId(email string) (string, error) {
	var result UserID
	var c Session
	s := UseDb(&c)
	err := s.CreateSession()
	if err != nil {
		return "", err
	}
	err = s.GetCollection("users").Find(bson.M{"email": email}).One(&result)
	if err != nil {
		return "", err
	}
	defer s.CloseSession()
	return string(result.ID), nil
}

func ValidateToken(w http.ResponseWriter, r *http.Request) (string, error) {
	authorizationHeader := r.Header.Get("authorization")
	if authorizationHeader != "" {
		token, error := jwt.Parse(authorizationHeader, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return []byte("secret"), nil
		})
		if error != nil {
			return "", errors.New("Invalid authorization token")
		}
		if token.Valid {
			c, err := r.Cookie(authorizationHeader)
			if err != nil {
				return "", errors.New("Tocken not found in cookie")
			}
			if id, err := GetUserId(c.Value); err == nil {
				return id, nil
			}
			return "", err
		}
		return "", errors.New("Timing is everything")
	}
	return "", errors.New("An authorization header is required")
}
