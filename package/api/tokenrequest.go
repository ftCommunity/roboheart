package api

type TokenRequest struct {
	Token string `json:"token" query:"token"`
}
