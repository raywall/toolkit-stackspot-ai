package auth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/raywall/toolkit-stackspot-ai/pkg/clients"
	"github.com/raywall/toolkit-stackspot-ai/pkg/types"
	"github.com/raywall/toolkit-stackspot-ai/services/auth"
)

func newTestClient(t *testing.T, handler http.Handler) *clients.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return clients.New(
		clients.WithBaseURL(srv.URL),
		clients.WithTimeout(5*time.Second),
	)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func TestAuthService_GenerateToken(t *testing.T) {
	ctx := context.Background()
	expected := auth.TokenResponse{
		AccessToken: "tok_abc123",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "oidc/oauth/token") || r.Method != http.MethodPost {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		writeJSON(w, http.StatusOK, expected)
	}))

	svc := auth.NewAuthService(client)

	resp, err := svc.GenerateToken(ctx, &types.Credentials{
		ClientID:     "test-id",
		ClientSecret: "test-secret",
	})
	if err != nil {
		t.Fatalf("esperado sem erro, obteve: %v", err)
	}
	if resp.AccessToken != expected.AccessToken {
		t.Errorf("token esperado %q, obteve %q", expected.AccessToken, resp.AccessToken)
	}
	if client.Token != expected.AccessToken {
		t.Error("token não foi salvo no client base")
	}
}

func TestAuthService_GenerateToken_Error(t *testing.T) {
	ctx := context.Background()
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"code":    "UNAUTHORIZED",
			"message": "credenciais inválidas",
		})
	}))

	svc := auth.NewAuthService(client)

	_, err := svc.GenerateToken(ctx, &types.Credentials{
		ClientID:     "wrong",
		ClientSecret: "wrong",
	})
	if err == nil {
		t.Fatal("esperado erro, obteve nil")
	}
	if !clients.IsUnauthorized(err) {
		t.Errorf("esperado erro Unauthorized, obteve: %v", err)
	}
}

func TestAuthService_RevokeToken(t *testing.T) {
	ctx := context.Background()
	revoked := false

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oidc/oauth/token/revoke":
			revoked = true
			w.WriteHeader(http.StatusNoContent)
		default:
			http.NotFound(w, r)
		}
	}))

	// Simulando um client já autenticado
	client.Token = "tok_xyz"
	svc := auth.NewAuthService(client)

	if err := svc.RevokeToken(ctx); err != nil {
		t.Fatalf("erro ao revogar token: %v", err)
	}
	if !revoked {
		t.Error("endpoint de revogação não foi chamado")
	}
	if client.Token != "" {
		t.Error("token não foi removido do client após revogação")
	}
}
