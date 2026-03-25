package auth

import "time"

// TokenResponse é retornado pelos endpoints de autenticação e renovação.
type TokenResponse struct {
	// AccessToken é o token de acesso Bearer.
	AccessToken string `json:"access_token"`
	// TokenType é o tipo do token (ex: "Bearer").
	TokenType string `json:"token_type"`
	// ExpiresIn é a duração de validade em segundos.
	ExpiresIn int `json:"expires_in"`
	// ExpiresAt é o timestamp calculado de expiração.
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	// RefreshExpiresIn é o duração de validade do refresh token em segundos
	RefreshExpiresIn int `json:"refresh_expires_in,omitempty"`
	// RefreshExpiresAt é o timestamp calculado de expiração.
	RefreshExpiresAt time.Time `json:"refresh_expires_at,omitempty"`
	// RefreshToken é o token para renovação (quando disponível).
	RefreshToken string `json:"refresh_token,omitempty"`
	// Scope lista os escopos concedidos.
	Scope string `json:"scope,omitempty"`
	// NotBeforePolicy
	NotBeforePolicy int `json:"not-before-policy,omitempty"`
	// SessionState
	SessionState string `json:"session_state,omitempty"`
}

// IsExpired retorna true se o token já expirou.
func (t *TokenResponse) IsExpired() bool {
	if t.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(t.ExpiresAt)
}
