package dto

type LangBinding struct {
	FromISO string `json:"from_iso" param:"from_iso" validate:"required"`
	ToISO   string `json:"to_iso" param:"to_iso" validate:"required"`
	Active  bool   `json:"active,omitempty"`
}
