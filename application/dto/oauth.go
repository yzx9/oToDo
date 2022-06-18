package dto

type OAuthRedirector struct {
	RedirectURI string `json:"redirectURI"`
}

type OAuthPayload struct {
	Code  string `json:"code"`
	State string `json:"state"`
}
