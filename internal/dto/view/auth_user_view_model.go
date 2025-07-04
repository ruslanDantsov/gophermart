package view

//go:generate easyjson -all auth_user_view_model.go
type AuthUserViewModel struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}
