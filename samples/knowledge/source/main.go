package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/raywall/toolkit-stackspot-ai/internal/auth"
	"github.com/raywall/toolkit-stackspot-ai/internal/knowledge/source"
	"github.com/raywall/toolkit-stackspot-ai/pkg/clients"
	"github.com/raywall/toolkit-stackspot-ai/pkg/types"
)

func main() {
	ctx := context.Background()

	client := clients.New(clients.WithBaseURL(os.Getenv("API_BASE_URL")))
	client.TokenProvider = func(ctx context.Context) error {
		_, err := auth.NewAuthService(client).GenerateToken(ctx, &types.Credentials{
			ClientID:     os.Getenv("API_CLIENT_ID"),
			ClientSecret: os.Getenv("API_CLIENT_SECRET"),
		})
		return err
	}

	ksSvc := source.NewKnowledgeSourceService(client)

	fmt.Println("→ Criando Knowledge Source...")
	ks, err := ksSvc.Create(ctx, &source.CreateKnowledgeSourceRequest{
		Name:        "Base de Documentação API",
		Description: "Fonte contendo especificações técnicas",
		Type:        types.KnowledgeSourceTypeDocument,
	})
	if err != nil {
		log.Fatalf("Erro ao criar: %v", err)
	}
	fmt.Printf("✓ Source criada: ID=%s, Nome=%s\n", ks.ID, ks.Name)

	fmt.Println("\n→ Listando Sources Ativas...")
	list, err := ksSvc.List(ctx, &source.ListKnowledgeSourcesParams{
		Status:     types.KnowledgeSourceStatusActive,
		Pagination: types.Pagination{PageSize: 5},
	})
	if err != nil {
		log.Fatalf("Erro ao listar: %v", err)
	}
	fmt.Printf("✓ %d sources encontradas na primeira página.\n", len(list.Items))

	fmt.Println("\n→ Removendo Source...")
	if err := ksSvc.Delete(ctx, ks.ID); err != nil {
		log.Fatalf("Erro ao deletar: %v", err)
	}
	fmt.Println("✓ Source removida com sucesso!")
}
