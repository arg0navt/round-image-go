package db

import (
	"errors"
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// URL path top db.
const URL = "localhost:27017"

// DB name of db
const DB = "rimg"

// The UserID return this struct if will find user with Id
type UserID struct {
	ID bson.ObjectId `json:"id" bson:"_id"`
}

// The UseDb this is interface for work with db
type UseDb interface {
	CreateSession() error
	CloseSession()
	GetCollection(name string) *mgo.Collection
	FindUserByID(id string, result interface{}) error
}

// Session struct of interface UseDb
type Session struct {
	Value *mgo.Session
}

// CreateSession It is create new session to db
func (s *Session) CreateSession() error {
	connect, err := mgo.Dial(URL)
	if err != nil {
		return err
	}
	s.Value = connect
	return nil
}

// CloseSession It is close of open session
func (s *Session) CloseSession() {
	s.Value.Close()
}

// GetCollection It is return collection with the necessary name
func (s *Session) GetCollection(name string) *mgo.Collection {
	return s.Value.DB(DB).C(name)
}

// FindUserByID  It is find user by her id
func (s *Session) FindUserByID(id string, result interface{}) error {
	err := s.GetCollection("users").FindId(bson.ObjectIdHex(id)).One(&result)
	if err != nil {
		return err
	}
	return nil
}

// ThereIsUserEmail test for email of valid
func ThereIsUserEmail(s UseDb, email string) bool {
	result, _ := s.GetCollection("users").Find(bson.M{"email": email}).Count()
	if result != 0 {
		return true
	}
	return false
}

// GetUserID it is find user and return UserID
func GetUserID(email string) (string, error) {
	var result UserID
	var s UseDb = &Session{}
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

// ValidateToken it is get token from header request and validate him.
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
			if id, err := GetUserID(c.Value); err == nil {
				return id, nil
			}
			return "", err
		}
		return "", errors.New("Timing is everything")
	}
	return "", errors.New("An authorization header is required")
}
