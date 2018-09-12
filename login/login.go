package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"
	"unicode/utf8"

	"../db"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fatih/structs"
	"gopkg.in/mgo.v2/bson"
)

// maxAge this is time for save of token in cookie
const maxAge = 86400

// RequestLogInSignUp  struct reques of /log_in /sign_up
type RequestLogInSignUp struct {
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
}

// User push this info to bd /sign_up
type User struct {
	FirstName      string `json:"first_name" bson:"first_name"`
	LastName       string `json:"last_name" bson:"last_name"`
	Email          string `json:"email" bson:"email"`
	Password       string `json:"password" bson:"password"`
	Verification   bool   `json:"verification" bson:"verification"`
	DateLastActive int64  `json:"dateLastActive" bson:"dateLastActive"`
	DetailInfo     `json:"detailInfo" bson:"detailInfo"`
}

// DetailInfo detail information by user
type DetailInfo struct {
	Avatar          string `json:"avatar" bson:"avatar"`
	ImageBackground string `json:"imageBackground" bson:"imageBackground"`
	StatusMessage   string `json:"statusMessage" bson:"statusMessage"`
}

// Exception return after /check_token /log_out
type Exception struct {
	Message string `json:"message"`
}

// JwtToken return this message /log_in /sign_up
type JwtToken struct {
	Token string `json:"token"`
}

// CreateUser create new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var target RequestLogInSignUp
	if ready := readyData(&target, w, r); ready == true {
		var s db.UseDb = &db.Session{}
		err := s.CreateSession()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer s.CloseSession()
		okValid, err := validateValuesSignUp(&target, s)
		if okValid == false {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		newUser := createBasicStruct(&target)
		err = s.GetCollection("users").Insert(&newUser)
		if err != nil {
			fmt.Println(err)
		}
		token := createToken(w, target.Email)
		json.NewEncoder(w).Encode(JwtToken{Token: token})
	}
}

// LogIn return new token for user
func LogIn(w http.ResponseWriter, r *http.Request) {
	var target RequestLogInSignUp
	if ready := readyData(&target, w, r); ready == true {
		okValid, err := validateValuesLogIn(&target)
		if okValid == false {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var s db.UseDb = &db.Session{}
		err = s.CreateSession()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer s.CloseSession()
		count, errCount := s.GetCollection("users").Find(bson.M{"email": target.Email, "password": target.Password}).Count()
		if errCount != nil {
			http.Error(w, errCount.Error(), http.StatusInternalServerError)
			return
		}
		if count == 0 {
			http.Error(w, "User not found", http.StatusInternalServerError)
			return
		}
		token := createToken(w, target.Email)
		json.NewEncoder(w).Encode(JwtToken{Token: token})
	}
}

// LogOut delete token from cookie
func LogOut(w http.ResponseWriter, r *http.Request) {
	token, err := db.ValidateToken(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	c := http.Cookie{
		Name:   token,
		MaxAge: -1,
	}
	http.SetCookie(w, &c)
	json.NewEncoder(w).Encode(Exception{Message: "ok"})
}

// CheckToken validate user token
func CheckToken(w http.ResponseWriter, r *http.Request) {
	_, err := db.ValidateToken(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(Exception{Message: "ok"})
}

func createBasicStruct(t *RequestLogInSignUp) User {
	return User{
		FirstName:      t.FirstName,
		LastName:       t.LastName,
		Email:          t.Email,
		Password:       t.Password,
		Verification:   false,
		DateLastActive: time.Now().Unix(),
	}
}

func readyData(target *RequestLogInSignUp, w http.ResponseWriter, r *http.Request) bool {
	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return false
	}
	err := json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return false
	}

	return true
}

func createToken(w http.ResponseWriter, e string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": e,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"iat":   time.Now().Unix(),
	})
	tokenString, _ := token.SignedString([]byte("secret"))
	http.SetCookie(w, &http.Cookie{
		Name:   tokenString,
		Value:  e,
		MaxAge: maxAge,
	})
	return tokenString
}

func validateValuesSignUp(values *RequestLogInSignUp, s db.UseDb) (bool, error) {
	mapValues := structs.Map(values)
	for key, value := range mapValues {
		switch key {
		case "FirstName", "LastName":
			nameV := validate(value.(string), `[a-zA-Z]`, 2, 32)
			if nameV == false {
				return false, errors.New("name error")
			}
		case "Email":
			emailV := validate(value.(string), `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`, 5, 50)
			if db.ThereIsUserEmail(s, value.(string)) == true {
				return false, errors.New("occupied email")
			}
			if emailV == false {
				return false, errors.New("email error")
			}
		case "Password":
			passV := validate(value.(string), `[a-zA-Z0-9]`, 8, 20)
			if passV == false {
				return false, errors.New("password error")
			}
		}
	}
	return true, nil
}

func validateValuesLogIn(values *RequestLogInSignUp) (bool, error) {
	mapValues := structs.Map(values)
	for key, value := range mapValues {
		switch key {
		case "Email":
			emailV := validate(value.(string), `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`, 5, 50)
			if emailV == false {
				return false, errors.New("email error")
			}
		case "Password":
			passV := validate(value.(string), `[a-zA-Z0-9]`, 8, 20)
			if passV == false {
				return false, errors.New("password error")
			}
		}
	}
	return true, nil
}

func validate(t string, reg string, min int, max int) bool {
	lenT := utf8.RuneCountInString(t)
	regE := regexp.MustCompile(reg)
	return regE.MatchString(t) && (lenT >= min && lenT <= max)
}
