package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/raywall/toolkit-stackspot-ai/pkg/clients"
	"github.com/raywall/toolkit-stackspot-ai/pkg/types"
)

// AuthService gerencia autenticação e geração de tokens.
type AuthService struct {
	client *clients.Client
}

// NewAuthService cria uma nova instância do serviço de autenticação.
func NewAuthService(c *clients.Client) *AuthService {
	return &AuthService{client: c}
}

// GenerateToken autentica as credenciais e retorna um TokenResponse.
// O token é automaticamente armazenado no Client base para uso nas demais requisições.
func (s *AuthService) GenerateToken(ctx context.Context, creds *types.Credentials) (*TokenResponse, error) {
	if creds.GrantType == "" {
		creds.GrantType = "client_credentials"
	}

	data := url.Values{}
	data.Set("grant_type", creds.GrantType)
	data.Set("client_id", creds.ClientID)
	data.Set("client_secret", creds.ClientSecret)

	encodedData := data.Encode()
	payload := strings.NewReader(encodedData)
	authUrl := fmt.Sprintf("%s/%s/oidc/oauth/token", os.Getenv("AUTH_BASE_URL"), os.Getenv("AUTH_REALM"))

	req, err := s.client.NewRequest(ctx, http.MethodPost, authUrl, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var resp TokenResponse
	if err := s.client.Do(req, &resp); err != nil {
		return nil, err
	}

	if resp.ExpiresIn > 0 {
		resp.ExpiresAt = time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second)
		resp.RefreshExpiresAt = time.Now().Add(time.Duration(resp.RefreshExpiresIn) * time.Second)
	}

	// Atualiza o token no cliente HTTP base
	s.client.Token = resp.AccessToken
	return &resp, nil
}

// RefreshToken renova o token de acesso usando um refresh token.
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	payload := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
	}

	authUrl := fmt.Sprintf("%s/%soidc/oauth/token/refresh", os.Getenv("AUTH_BASE_URL"), os.Getenv("AUTH_REALM"))

	req, err := s.client.NewRequest(ctx, http.MethodPost, authUrl, payload)
	if err != nil {
		return nil, err
	}

	var resp TokenResponse
	if err := s.client.Do(req, &resp); err != nil {
		return nil, err
	}

	if resp.ExpiresIn > 0 {
		resp.ExpiresAt = time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second)
		resp.RefreshExpiresAt = time.Now().Add(time.Duration(resp.RefreshExpiresIn) * time.Second)
	}

	// Atualiza o token no cliente HTTP base
	s.client.Token = resp.AccessToken
	return &resp, nil
}

// RevokeToken invalida o token de acesso atual.
func (s *AuthService) RevokeToken(ctx context.Context) error {
	if s.client.Token == "" {
		return nil
	}

	payload := map[string]string{"token": s.client.Token}
	authUrl := fmt.Sprintf("%s/%soidc/oauth/token/revoke", os.Getenv("AUTH_BASE_URL"), os.Getenv("AUTH_REALM"))

	req, err := s.client.NewRequest(ctx, http.MethodPost, authUrl, payload)
	if err != nil {
		return err
	}

	if err := s.client.Do(req, nil); err != nil {
		return err
	}

	// Remove o token do cliente HTTP base após revogação
	s.client.Token = ""
	return nil
}
