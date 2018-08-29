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

type RequestLogInSignUp struct {
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
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
		err := db.GetUsers().Insert(target)
		if err != nil {
			fmt.Println(err)
		}
		token, error := createToken(target.Email)
		if error != nil {
			fmt.Println(error)
		}
		json.NewEncoder(w).Encode(JwtToken{Token: token})
	}
}

func LogIn(w http.ResponseWriter, r *http.Request) {
	var target RequestLogInSignUp
	if ready := readyData(&target, w, r); ready == true {
		token, error := createToken(target.Email)
		if error != nil {
			fmt.Println(error)
		}
		json.NewEncoder(w).Encode(JwtToken{Token: token})
	}
}

func CheckToken(w http.ResponseWriter, r *http.Request) {
	if v := validateToken(w, r); v == true {
		json.NewEncoder(w).Encode(Exception{Message: "ok"})
	}
}

func validateToken(w http.ResponseWriter, r *http.Request) bool {
	authorizationHeader := r.Header.Get("authorization")
	if authorizationHeader != "" {
		token, error := jwt.Parse(authorizationHeader, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return []byte("secret"), nil
		})
		if error != nil {
			http.Error(w, "Invalid authorization token", 400)
			return false
		}
		if token.Valid {
			return true
		} else {
			http.Error(w, "Invalid authorization token", 400)
		}

	} else {
		http.Error(w, "An authorization header is required", 400)
	}
	return false
}

func readyData(target *RequestLogInSignUp, w http.ResponseWriter, r *http.Request) bool {
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return false
	}
	err := json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return false
	}
	okValid, text := validateValues(target)
	if okValid == false {
		http.Error(w, text, 400)
		return false
	}
	return true
}

func createToken(e string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": e,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"iat":   time.Now().Unix(),
	})
	return token.SignedString([]byte("secret"))
}

func validateValues(values *RequestLogInSignUp) (bool, string) {
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
			if db.ThereIsUser(value.(string)) == true {
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

func validate(t string, reg string, min int, max int) bool {
	lenT := utf8.RuneCountInString(t)
	regE := regexp.MustCompile(reg)
	return regE.MatchString(t) && (lenT >= min && lenT <= max)
}
