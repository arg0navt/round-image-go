package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"
	"unicode/utf8"

	"../db"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fatih/structs"
)

const maxAge = 86400 // duration valid token

type RequestLogInSignUp struct {
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
}

type User struct {
	FirstName      string `json:"first_name" bson:"first_name"`
	LastName       string `json:"last_name" bson:"last_name"`
	Email          string `json:"email" bson:"email"`
	Password       string `json:"password" bson:"password"`
	Verification   bool   `json:"verification" bson:"verification"`
	DateLastActive int64  `json:"dateLastActive" bson:"dateLastActive"`
	DetailInfo     `json:"detailInfo" bson:"detailInfo"`
}

type Img struct {
	Name string `json:"name" bson:"name"`
	Url  string `json:"url" bson:"url"`
}

type DetailInfo struct {
	Avatar          string `json:"avatar" bson:"avatar"`
	ImageBackground string `json:"imageBackground" bson:"imageBackground"`
	StatusMessage   string `json:"statusMessage" bson:"statusMessage"`
}

type Exception struct {
	Message string `json:"message"`
}

type JwtToken struct {
	Token string `json:"token"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var target RequestLogInSignUp
	if ready := readyData(&target, w, r); ready == true {
		okValid, text := validateValuesSignUp(&target)
		if okValid == false {
			http.Error(w, text, http.StatusInternalServerError)
			return
		}
		newUser := createBasicStruct(&target)
		err := db.GetCollection("users").Insert(&newUser)
		if err != nil {
			fmt.Println(err)
		}
		token := createToken(w, target.Email)
		json.NewEncoder(w).Encode(JwtToken{Token: token})
	}
}

func LogIn(w http.ResponseWriter, r *http.Request) {
	var target RequestLogInSignUp
	if ready := readyData(&target, w, r); ready == true {
		okValid, text := validateValuesLogIn(&target)
		if okValid == false {
			http.Error(w, text, http.StatusInternalServerError)
			return
		}
		token := createToken(w, target.Email)
		json.NewEncoder(w).Encode(JwtToken{Token: token})
	}
}

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

func validateValuesSignUp(values *RequestLogInSignUp) (bool, string) {
	mapValues := structs.Map(values)
	for key, value := range mapValues {
		switch key {
		case "FirstName", "LastName":
			nameV := validate(value.(string), `[a-zA-Z]`, 2, 32)
			if nameV == false {
				return false, "name error"
			}
		case "Email":
			emailV := validate(value.(string), `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`, 5, 50)
			if db.ThereIsUserEmail(value.(string)) == true {
				return false, "occupied email"
			}
			if emailV == false {
				return false, "email error"
			}
		case "Password":
			passV := validate(value.(string), `[a-zA-Z0-9]`, 8, 20)
			if passV == false {
				return false, "password error"
			}
		}
	}
	return true, "validate"
}

func validateValuesLogIn(values *RequestLogInSignUp) (bool, string) {
	mapValues := structs.Map(values)
	for key, value := range mapValues {
		switch key {
		case "Email":
			emailV := validate(value.(string), `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`, 5, 50)
			if emailV == false {
				return false, "email error"
			}
		case "Password":
			passV := validate(value.(string), `[a-zA-Z0-9]`, 8, 20)
			if passV == false {
				return false, "password error"
			}
		}
	}
	return true, "validate"
}

func validate(t string, reg string, min int, max int) bool {
	lenT := utf8.RuneCountInString(t)
	regE := regexp.MustCompile(reg)
	return regE.MatchString(t) && (lenT >= min && lenT <= max)
}
