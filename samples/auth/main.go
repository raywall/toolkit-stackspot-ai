package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/raywall/toolkit-stackspot-ai/internal/auth"
	"github.com/raywall/toolkit-stackspot-ai/pkg/clients"
	"github.com/raywall/toolkit-stackspot-ai/pkg/types"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := clients.New(clients.WithBaseURL(os.Getenv("API_BASE_URL")))
	authSvc := auth.NewAuthService(client)

	creds := &types.Credentials{
		ClientID:     os.Getenv("API_CLIENT_ID"),
		ClientSecret: os.Getenv("API_CLIENT_SECRET"),
	}

	fmt.Println("→ Autenticando...")
	token, err := authSvc.GenerateToken(ctx, creds)
	if err != nil {
		log.Fatalf("Erro ao autenticar: %v", err)
	}

	fmt.Printf("✓ Token gerado com sucesso!\n")
	fmt.Printf("  Access Token: %s...\n", token.AccessToken[:15]) // Mostra apenas o início
	fmt.Printf("  Expira em: %d segundos\n", token.ExpiresIn)

	fmt.Println("\n→ Revogando token...")
	if err := authSvc.RevokeToken(ctx); err != nil {
		log.Fatalf("Erro ao revogar token: %v", err)
	}
	fmt.Println("✓ Token revogado com sucesso!")
}
