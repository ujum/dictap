package dto

type User struct {
	UID   string `json:"uid"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type UserUpdate struct {
	UID   string `json:"-" swaggerignore:"true"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type UserCreate struct {
	User
	Password  string `json:"password"`
	OAuthUser bool   `json:"-" swaggerignore:"true"`
}

type UserCredentials struct {
	Email    string `json:"email,required"`
	Password string ` json:"password,required"`
}
