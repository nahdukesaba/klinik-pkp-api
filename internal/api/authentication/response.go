package authentication

// AUTH RESPONSE
type AddAuthenticationResponse struct {
	ID           string `json:"id"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// REFRESH RESPONSE CONTAINS NEW ACCESS TOKEN ONLY
type RefreshAuthenticationResponse struct {
	AccessToken string `json:"accessToken"`
}
