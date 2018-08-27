package user

import (
	"encoding/json"
	"net/http"
	"regexp"
	"unicode/utf8"

	"github.com/fatih/structs"
)

// type User struct {
// 	ID    int    `json:"value" bson:"_id,omitempty"`
// 	Email string `json:"email"`
// }

type RequestCreateUser interface {
	ValidateValues() (bool, string)
}

type NewUser struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password"`
}

func (u NewUser) ValidateValues() (bool, string) {
	mapUser := structs.Map(u)
	for key, value := range mapUser {
		switch key {
		case "FirstName", "LastName":
			nameV := validate(value.(string), `[a-zA-Z]`, 2, 32)
			if nameV == false {
				return false, "name error "
			}
		case "Email":
			emailV := validate(value.(string), `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`, 5, 50)
			if emailV == false {
				return false, "email error "
			}
		case "Password":
			if value.(string) != mapUser["RepeatPassword"].(string) {
				return false, "repeat password != password "
			}
			passV := validate(value.(string), `[a-zA-Z0-9]`, 8, 20)
			if passV == false {
				return false, "password error "
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

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var target NewUser
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	var data RequestCreateUser = target
	okValid, text := data.ValidateValues()
	if okValid == false {
		http.Error(w, text, 400)
		return
	}
}
