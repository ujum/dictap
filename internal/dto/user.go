package dto

type User struct {
	UID         string        `json:"uid"`
	Name        string        `json:"name"`
	Email       string        `json:"email"`
	Phone       string        `json:"phone"`
	OAuthUser   bool          `json:"oauth_user"`
	LangBinding []LangBinding `json:"lang_binding,omitempty"`
}

type UserUpdate struct {
	UID         string        `json:"-" swaggerignore:"true"`
	Name        string        `json:"name"`
	Phone       string        `json:"phone"`
	LangBinding []LangBinding `json:"lang_binding,omitempty"`
}

type UserCreate struct {
	Name        string              `json:"name"`
	Email       string              `json:"email"`
	Phone       string              `json:"phone"`
	Password    string              `json:"password"`
	OAuthUser   bool                `json:"-" swaggerignore:"true"`
	LangBinding []LangBindingCreate `json:"lang_binding"`
}

type UserCredentials struct {
	Email    string `json:"email,required"`
	Password string ` json:"password,required"`
}

type ChangeUserPassword struct {
	OldPassword string ` json:"old_password,required"`
	Password    string ` json:"password,required"`
}

type SetUserPassword struct {
	Password string ` json:"password,required"`
}
