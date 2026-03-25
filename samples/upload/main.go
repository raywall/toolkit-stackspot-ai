package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/raywall/toolkit-stackspot-ai/pkg/clients"
	"github.com/raywall/toolkit-stackspot-ai/pkg/types"
	"github.com/raywall/toolkit-stackspot-ai/services/auth"
	"github.com/raywall/toolkit-stackspot-ai/services/knowledge/object"
)

func main() {
	ctx := context.Background()

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

	// 1. Criar um arquivo temporário de exemplo
	dir := os.TempDir()
	fileName := "tabela_precos.csv"
	filePath := filepath.Join(dir, fileName)

	err := os.WriteFile(filePath, []byte("produto,preco\nCelular,1500\nNotebook,4500"), 0644)
	if err != nil {
		log.Fatalf("Erro ao criar arquivo temp: %v", err)
	}
	defer os.Remove(filePath) // Limpa o arquivo temp no final

	// 2. Gerar Intenção de Upload
	fmt.Printf("→ Gerando intenção de upload para '%s'...\n", fileName)
	intent, err := koSvc.GenerateUpload(ctx, sourceID, &object.CreateUploadRequest{
		FileName: fileName,
		MimeType: "text/csv",
	})
	if err != nil {
		log.Fatalf("Erro ao gerar upload: %v", err)
	}
	fmt.Printf("✓ Intenção gerada! UploadID: %s\n", intent.UploadID)

	// 3. Fazer o Upload físico
	fmt.Println("\n→ Fazendo upload do arquivo (multipart/form-data)...")
	err = koSvc.UploadFile(ctx, sourceID, intent.UploadID, dir, fileName)
	if err != nil {
		log.Fatalf("Erro ao enviar arquivo: %v", err)
	}
	fmt.Println("✓ Upload concluído com sucesso!")
}
