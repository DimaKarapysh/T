package redis

type SessionData struct {
	RefreshTokenHash string `json:"token"`
	UserAgent        string `json:"user_agent"`
	IP               string `json:"ip"`
}
