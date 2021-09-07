package domain

import (
	"time"
)

type User struct {
	ID           string    `bson:"_id,omitempty"`
	UID          string    `bson:"uid"`
	Name         string    `bson:"name"`
	Email        string    `bson:"email,omitempty"`
	Phone        string    `bson:"phone"`
	Password     string    `bson:"password,omitempty"`
	RegisteredAt time.Time `bson:"registeredAt"`
	OAuthUser    bool      `bson:"oauth_user"`
}
