package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/raywall/toolkit-stackspot-ai/internal/auth"
	"github.com/raywall/toolkit-stackspot-ai/internal/knowledge/object"
	"github.com/raywall/toolkit-stackspot-ai/pkg/clients"
	"github.com/raywall/toolkit-stackspot-ai/pkg/types"
)

func main() {
	ctx := context.Background()

	// Requer que você passe um SOURCE_ID existente via variável de ambiente para testar
	sourceID := os.Getenv("SOURCE_ID")
	if sourceID == "" {
		log.Fatal("A variável SOURCE_ID é obrigatória para este exemplo")
	}

	client := clients.New(clients.WithBaseURL(os.Getenv("API_BASE_URL")))
	client.TokenProvider = func(ctx context.Context) error {
		_, err := auth.NewAuthService(client).GenerateToken(ctx, &types.Credentials{
			ClientID:     os.Getenv("API_CLIENT_ID"),
			ClientSecret: os.Getenv("API_CLIENT_SECRET"),
		})
		return err
	}

	koSvc := object.NewKnowledgeObjectService(client)

	fmt.Println("→ Criando Knowledge Object (Texto)...")
	ko, err := koSvc.Create(ctx, sourceID, &object.CreateKnowledgeObjectRequest{
		Title:   "Políticas da Empresa",
		Content: "O horário de trabalho padrão é das 09h às 18h.",
		Tags:    []string{"rh", "politicas"},
	})
	if err != nil {
		log.Fatalf("Erro ao criar objeto: %v", err)
	}
	fmt.Printf("✓ Objeto criado: ID=%s, Título=%s\n", ko.ID, ko.Title)

	fmt.Println("\n→ Deletando todos os objetos desta Source...")
	if err := koSvc.DeleteAll(ctx, sourceID); err != nil {
		log.Fatalf("Erro ao limpar objetos: %v", err)
	}
	fmt.Println("✓ Todos os objetos removidos com sucesso!")
}
