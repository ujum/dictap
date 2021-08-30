package dto

type User struct {
	Uid   string `json:"uid"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type UserUpdate struct {
	Uid   string `json:"-" swaggerignore:"true"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type UserCreate struct {
	User
	Password string `json:"password"`
}
