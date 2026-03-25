package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/raywall/toolkit-stackspot-ai/internal/agent"
	"github.com/raywall/toolkit-stackspot-ai/internal/auth"
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

	agentSvc := agent.NewAgentService(client)

	fmt.Println("→ Criando Agente...")
	ag, err := agentSvc.Create(ctx, &agent.CreateAgentRequest{
		Name:         "Oráculo Tech",
		Model:        types.AgentModelGPT4o,
		SystemPrompt: "Você é um especialista em Go e arquitetura de software.",
		Temperature:  0.2,
	})
	if err != nil {
		log.Fatalf("Erro ao criar agente: %v", err)
	}
	fmt.Printf("✓ Agente criado: ID=%s\n", ag.ID)

	fmt.Println("\n→ Executando Agente (Síncrono)...")
	syncResp, err := agentSvc.Execute(ctx, ag.ID, &agent.ExecuteAgentRequest{
		Input: "Quais os benefícios de usar Go para microsserviços?",
	})
	if err != nil {
		log.Fatalf("Erro na execução síncrona: %v", err)
	}
	fmt.Printf("✓ Resposta: %s\n", syncResp.Output)
	if syncResp.Usage != nil {
		fmt.Printf("  Tokens usados: %d\n", syncResp.Usage.TotalTokens)
	}

	fmt.Println("\n→ Executando Agente (Streaming)...")
	fmt.Print("✓ Resposta: ")
	err = agentSvc.ExecuteStream(ctx, ag.ID, &agent.ExecuteAgentRequest{
		Input: "Resuma em uma frase o que é Kubernetes.",
	}, func(event agent.StreamEvent) error {
		if event.Type == "delta" {
			fmt.Print(event.Delta)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("\nErro no streaming: %v", err)
	}

	fmt.Println("→ Removendo Agente...")
	if err := agentSvc.Delete(ctx, ag.ID); err != nil {
		log.Fatalf("Erro ao deletar: %v", err)
	}
	fmt.Println("✓ Agente removido com sucesso!")
}
