package models

import (
	"encoding/json"
	"time"
)

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`

	ExpirationTime int `json:"-"`
}

func (a *AccessTokenResponse) UnmarshalJSON(data []byte) error {
	type Alias AccessTokenResponse
	aux := &struct {
		*Alias
		ExpiresIn      int `json:"expires_in"`
		ExpirationTime int `json:"-"` // exclude from default decoding
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	a.ExpirationTime = int(time.Now().Unix()) + aux.ExpiresIn

	return nil
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`

	ExpirationTime int `json:"-"`
}

func (r *RefreshTokenResponse) UnmarshalJSON(data []byte) error {
	type Alias RefreshTokenResponse
	aux := &struct {
		*Alias
		ExpiresIn      int `json:"expires_in"`
		ExpirationTime int `json:"-"` // exclude from default decoding
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	r.ExpirationTime = int(time.Now().Unix()) + aux.ExpiresIn

	return nil
}
