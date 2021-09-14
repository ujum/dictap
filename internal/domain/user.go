package domain

import (
	"time"
)

type LangBinding struct {
	LangFromISO string `bson:"lang_from_iso"`
	LangToISO   string `bson:"lang_to_iso"`
	Active      bool   `bson:"active,omitempty"`
}

type User struct {
	ID           string        `bson:"_id,omitempty"`
	UID          string        `bson:"uid"`
	Name         string        `bson:"name"`
	Email        string        `bson:"email,omitempty"`
	Phone        string        `bson:"phone"`
	Password     string        `bson:"password,omitempty"`
	RegisteredAt time.Time     `bson:"registeredAt"`
	LangBinding  []LangBinding `bson:"lang_binding"`
	OAuthUser    bool          `bson:"oauth_user"`
}
