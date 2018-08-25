package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fatih/structs"
)

// type User struct {
// 	ID    int    `json:"value" bson:"_id,omitempty"`
// 	Email string `json:"email"`
// }

type RequestCreateUser interface {
	ValidateEmptyValues() bool
}

type NewUser struct {
	FirstName      string `json:"first_name" validate:"presence,min=2,max=32"`
	LastName       string `json:"last_name" validate:"presence,min=2,max=32"`
	Email          string `json:"email" validate:"email,required"`
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password"`
}

func (u NewUser) ValidateEmptyValues() bool {
	n := structs.Values(u)
	for _, i := range n {
		if i == "" {
			return false
		}
	}
	return true
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
	fmt.Println(data.ValidateEmptyValues())
}
