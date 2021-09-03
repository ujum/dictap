package dto

type Word struct {
	Name    string `json:"name"`
	AddedAt string `json:"added_at"`
}

type WordCreate struct {
	Name    string `json:"name"`
	GroupId string `json:"group_id"`
}

type WordGroupMovement struct {
	FromGroupId string `json:"from_group_id"`
	ToGroupId   string `json:"to_group_id"`
}

type WordGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Lang string `json:"lang_id"`
}

type WordGroupCreate struct {
	Name   string `json:"name"`
	Lang   string `json:"lang_id"`
	UserID string `json:"-" swaggerignore:"true"`
}
