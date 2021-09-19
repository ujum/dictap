package domain

import (
	"time"
)

type LangBinding struct {
	FromISO string `bson:"from_iso"`
	ToISO   string `bson:"to_iso"`
	Active  bool   `bson:"active,omitempty"`
}

type User struct {
	ID           string        `bson:"_id,omitempty"`
	UID          string        `bson:"uid,omitempty"`
	Name         string        `bson:"name,omitempty"`
	Email        string        `bson:"email,omitempty"`
	Phone        string        `bson:"phone,omitempty"`
	Password     string        `bson:"password,omitempty"`
	RegisteredAt time.Time     `bson:"registeredAt,omitempty"`
	LangBinding  []LangBinding `bson:"lang_binding,omitempty"`
	OAuthUser    bool          `bson:"oauth_user,omitempty"`
}
