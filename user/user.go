package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"
	"unicode/utf8"

	"gopkg.in/mgo.v2/bson"

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
	FirstName      string  `json:"first_name" bson:"first_name"`
	LastName       string  `json:"last_name" bson:"last_name"`
	Email          string  `json:"email" bson:"email"`
	Password       string  `json:"password" bson:"password"`
	Verification   bool    `json:"verification" bson:"verification"`
	DateLastActive int64   `json:"dateLastActive" bson:"dateLastActive"`
	Albums         []Album `json:"albums" bson:"albums"`
	DetailInfo     `json:"detailInfo" bson:"detailInfo"`
}

type Album struct {
	ID           bson.ObjectId `json:"id" bson:"_id"`
	Name         string        `json:"name" bson:"name"`
	TimeToCreate int64         `json:"timeToCreate" bson:"timeToCreate"`
	Description  string        `json:"description" bson:"description"`
	Images       []Img         `json:"images" bson:"images"`
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
			http.Error(w, text, 400)
			return
		}
		newUser := createBasicStruct(&target)
		err := db.GetUsers().Insert(&newUser)
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
			http.Error(w, text, 400)
			return
		}
		token := createToken(w, target.Email)
		json.NewEncoder(w).Encode(JwtToken{Token: token})
	}
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	if v, token := validateToken(w, r); v == true {
		c := http.Cookie{
			Name:   token,
			MaxAge: -1,
		}
		http.SetCookie(w, &c)
		json.NewEncoder(w).Encode(Exception{Message: "ok"})
	}
}

func CheckToken(w http.ResponseWriter, r *http.Request) {
	if v, _ := validateToken(w, r); v == true {
		json.NewEncoder(w).Encode(Exception{Message: "ok"})
	}
}

func createBasicStruct(t *RequestLogInSignUp) User {
	var defaultAlbum []Album
	emptyImages := make([]Img, 0)
	defaultAlbum = append(defaultAlbum, Album{
		ID:           bson.NewObjectId(),
		Name:         "default album",
		TimeToCreate: time.Now().Unix(),
		Images:       emptyImages,
	})
	return User{
		FirstName:      t.FirstName,
		LastName:       t.LastName,
		Email:          t.Email,
		Password:       t.Password,
		Verification:   false,
		DateLastActive: time.Now().Unix(),
		Albums:         defaultAlbum,
	}
}

func validateToken(w http.ResponseWriter, r *http.Request) (bool, string) {
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
			return false, "Invalid authorization token"
		}
		if token.Valid {
			_, err := r.Cookie(authorizationHeader)
			if err != nil {
				http.Error(w, "Tocken not found in cookie", 400)
				return false, "Tocken not found in cookie"
			}
			return true, authorizationHeader
		}
		http.Error(w, "Timing is everything", 400)
		return false, "Timing is everything"
	}
	http.Error(w, "An authorization header is required", 400)
	return false, "An authorization header is required"
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
		Value:  tokenString,
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
