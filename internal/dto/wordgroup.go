package dto

type WordGroup struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	LangBinding LangBinding `json:"lang_binding"`
}

type WordGroupCreate struct {
	Name        string      `json:"name"`
	UserUID     string      `json:"-" swaggerignore:"true"`
	LangBinding LangBinding `json:"lang_binding"`
}
