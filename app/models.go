package app

type AuthRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type AuthResponse struct {
    AuthToken      string `json:"authToken"`
    ReauthToken    string `json:"reauthToken"`
    ExpirationTime string `json:"expirationTime"`
}
