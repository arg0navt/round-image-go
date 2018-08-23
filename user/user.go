package user

type UserControl interface {
	CreateUser()
}

type User struct {
	ID    int    `json:"value" bson:"_id,omitempty"`
	Email string `json:"email"`
}
