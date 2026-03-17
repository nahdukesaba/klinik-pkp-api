package authentication

// AUTH RESPONSE
type AddAuthenticationV1Response struct {
	ID           string `json:"id"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"-"`
}

// REFRESH RESPONSE CONTAINS NEW ACCESS TOKEN ONLY
type RefreshAuthenticationV1Response struct {
	AccessToken string `json:"accessToken"`
}
