package dto

type User struct {
	UID         string        `json:"uid"`
	Name        string        `json:"name"`
	Email       string        `json:"email"`
	Phone       string        `json:"phone"`
	OAuthUser   bool          `json:"oauth_user"`
	LangBinding []LangBinding `json:"lang_binding,omitempty"`
	Roles       []string      `json:"roles"`
}

type UserUpdate struct {
	UID         string        `json:"-" swaggerignore:"true"`
	Name        string        `json:"name,omitempty"`
	Phone       string        `json:"phone,omitempty"`
	LangBinding []LangBinding `json:"lang_binding,omitempty" validate:"dive,required"`
}

type UserCreate struct {
	Name        string        `json:"name" validate:"required"`
	Email       string        `json:"email" validate:"required,email"`
	Phone       string        `json:"phone"`
	Password    string        `json:"password"`
	OAuthUser   bool          `json:"-" swaggerignore:"true"`
	LangBinding []LangBinding `json:"lang_binding" validate:"required,gt=0,dive,required"`
}

type UserCredentials struct {
	Email    string `json:"email" validate:"required"`
	Password string ` json:"password" validate:"required"`
}

type ChangeUserPassword struct {
	OldPassword string ` json:"old_password" validate:"required"`
	Password    string ` json:"password" validate:"required"`
}

type SetUserPassword struct {
	Password string ` json:"password" validate:"required"`
}
