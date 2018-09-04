package db

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
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

func ValidateToken(w http.ResponseWriter, r *http.Request) string {
	authorizationHeader := r.Header.Get("authorization")
	if authorizationHeader != "" {
		token, error := jwt.Parse(authorizationHeader, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return []byte("secret"), nil
		})
		if error != nil {
			http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
			return ""
		}
		if token.Valid {
			c, err := r.Cookie(authorizationHeader)
			if err != nil {
				http.Error(w, "Tocken not found in cookie", http.StatusUnauthorized)
				return ""
			}
			if id := GetUserId(c.Value); id != "" {
				return id
			}
			http.Error(w, "user not found", http.StatusUnauthorized)
			return ""
		}
		http.Error(w, "Timing is everything", http.StatusUnauthorized)
		return ""
	}
	http.Error(w, "An authorization header is required", http.StatusUnauthorized)
	return ""
}
