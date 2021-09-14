package dto

type LangBindingCreate struct {
	LangFromISO string `json:"lang_from_iso"`
	LangToISO   string `json:"lang_to_iso"`
}

type LangBinding struct {
	LangFromISO string `json:"lang_from_iso"`
	LangToISO   string `json:"lang_to_iso"`
	Active      bool   `json:"active,omitempty"`
}
