package dto

type WordGroup struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	UserUID     string      `json:"-" swaggerignore:"true,required"`
	LangBinding LangBinding `json:"lang_binding" validate:"required,dive,required"`
	Default     bool        `json:"default"`
}

type WordGroupCreate struct {
	Name        string      `json:"name"`
	UserUID     string      `json:"-" swaggerignore:"true,required"`
	LangBinding LangBinding `json:"lang_binding" validate:"required,dive,required"`
	Default     bool        `json:"-" swaggerignore:"true"`
}
