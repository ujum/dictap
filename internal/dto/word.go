package dto

type Word struct {
	Name    string `json:"name"`
	AddedAt string `json:"added_at"`
}

type WordCreate struct {
	Name    string `json:"name"`
	GroupID string `json:"group_id"`
}

type WordGroupMovement struct {
	FromGroupID string `json:"from_group_id"`
	ToGroupID   string `json:"to_group_id"`
}

type WordGroup struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	LangISO string `json:"lang_iso"`
}

type WordGroupCreate struct {
	Name    string `json:"name"`
	LangISO string `json:"lang_iso"`
	UserUID string `json:"-" swaggerignore:"true"`
}
