package dto

import "encoding/json"

type TokenDTO struct {
	AccessToken  json.RawMessage `json:"access_token,required" swaggertype:"string"`
	RefreshToken json.RawMessage `json:"refresh_token,required" swaggertype:"string"`
}
