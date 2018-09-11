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

func GetCollection(name string) *mgo.Collection {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	return session.DB(DB).C(name)
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
			if id := GetUserId(c.Value); id != "" {
				return id, nil
			}
			return "", errors.New("User not found")
		}
		return "", errors.New("Timing is everything")
	}
	return "", errors.New("An authorization header is required")
}
