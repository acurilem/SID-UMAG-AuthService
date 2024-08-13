package smaps

import "github.com/citiaps/SID-UMAG-AuthService/models"

type LoginModel struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refreshToken"`
	User         models.User `json:"user"`
}
