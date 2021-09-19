package dto

import "time"

type Word struct {
	Name    string    `json:"name"`
	AddedAt time.Time `json:"added_at"`
}

type WordCreate struct {
	Name    string `json:"name" validate:"required"`
	GroupID string `json:"group_id" validate:"required"`
}

type WordGroupMovement struct {
	FromGroupID string `json:"from_group_id" validate:"required"`
	ToGroupID   string `json:"to_group_id" validate:"required"`
}
