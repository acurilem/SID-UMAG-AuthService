package forms

type RefreshTokenForm struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
