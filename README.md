# agentsdk/go

Biblioteca Go para integração com a plataforma de Agentes. Fornece suporte completo a **Autenticação**, **Knowledge Sources**, **Knowledge Objects** e **Agentes** (incluindo execução síncrona e streaming SSE).

---

## Instalação

```bash
go get github.com/agentsdk/go
```

---

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    agentsdk "github.com/agentsdk/go"
)

func main() {
    ctx := context.Background()

    // Cria o client (autentica automaticamente na primeira requisição)
    client := agentsdk.New(
        "meu-client-id",
        "meu-client-secret",
        agentsdk.WithBaseURL("https://api.minha-plataforma.io/v1"),
    )

    // Executa um agente
    result, err := client.Agent.Execute(ctx, "agent-123", &agentsdk.ExecuteAgentRequest{
        Input: "Quais são os planos disponíveis?",
    })
    if err != nil {
        panic(err)
    }
    fmt.Println(result.Output)
}
```

---

## Funcionalidades

### 🔐 Autenticação

```go
// Gerar token (OAuth2 client_credentials)
token, err := client.Auth.GenerateToken(ctx, &agentsdk.Credentials{
    ClientID:     "meu-client-id",
    ClientSecret: "meu-client-secret",
})

// Renovar token com refresh token
newToken, err := client.Auth.RefreshToken(ctx, "meu-refresh-token")

// Revogar token
err := client.Auth.RevokeToken(ctx)

// Injetar token manualmente (sem autenticação automática)
client := agentsdk.New("", "", agentsdk.WithToken("meu-token"))
```

---

### 📚 Knowledge Source

```go
// Listar
sources, err := client.KnowledgeSource.List(ctx, &agentsdk.ListKnowledgeSourcesParams{
    Pagination: agentsdk.Pagination{Page: 1, PageSize: 20},
    Status:     agentsdk.KnowledgeSourceStatusActive,
    Type:       agentsdk.KnowledgeSourceTypeDocument,
    Search:     "documentação",
})

// Recuperar
source, err := client.KnowledgeSource.Get(ctx, "ks-123")

// Criar
source, err := client.KnowledgeSource.Create(ctx, &agentsdk.CreateKnowledgeSourceRequest{
    Name: "Documentação Técnica",
    Type: agentsdk.KnowledgeSourceTypeDocument,
    Tags: []string{"docs", "v2"},
})

// Atualizar (campos nulos são ignorados)
updated, err := client.KnowledgeSource.Update(ctx, "ks-123", &agentsdk.UpdateKnowledgeSourceRequest{
    Name:   agentsdk.StringPtr("Novo Nome"),
    Status: agentsdk.KnowledgeSourceStatusPtr(agentsdk.KnowledgeSourceStatusInactive),
})

// Deletar
err := client.KnowledgeSource.Delete(ctx, "ks-123")
```

---

### 📄 Knowledge Object

```go
// Listar (dentro de uma Knowledge Source)
objects, err := client.KnowledgeObject.List(ctx, "ks-123", &agentsdk.ListKnowledgeObjectsParams{
    Pagination: agentsdk.Pagination{Page: 1, PageSize: 50},
    Status:     agentsdk.KnowledgeObjectStatusActive,
})

// Recuperar
obj, err := client.KnowledgeObject.Get(ctx, "ks-123", "ko-456")

// Criar (via conteúdo direto)
obj, err := client.KnowledgeObject.Create(ctx, "ks-123", &agentsdk.CreateKnowledgeObjectRequest{
    Title:   "Guia de Instalação",
    Content: "Passo 1: baixe o instalador...",
})

// Criar (via URL)
obj, err := client.KnowledgeObject.Create(ctx, "ks-123", &agentsdk.CreateKnowledgeObjectRequest{
    Title:      "Manual do Usuário",
    ContentURL: "https://docs.exemplo.com/manual.pdf",
    MimeType:   "application/pdf",
})

// Atualizar
updated, err := client.KnowledgeObject.Update(ctx, "ks-123", "ko-456", &agentsdk.UpdateKnowledgeObjectRequest{
    Content: agentsdk.StringPtr("Conteúdo revisado..."),
})

// Deletar
err := client.KnowledgeObject.Delete(ctx, "ks-123", "ko-456")
```

---

### 🤖 Agentes

```go
// Listar
agents, err := client.Agent.List(ctx, &agentsdk.ListAgentsParams{
    Status: agentsdk.AgentStatusActive,
    Model:  agentsdk.AgentModelGPT4o,
})

// Recuperar
agent, err := client.Agent.Get(ctx, "agent-123")

// Criar
agent, err := client.Agent.Create(ctx, &agentsdk.CreateAgentRequest{
    Name:               "Assistente de Suporte",
    Model:              agentsdk.AgentModelGPT4o,
    SystemPrompt:       "Você é um assistente técnico especializado.",
    Temperature:        0.3,
    MaxTokens:          2048,
    KnowledgeSourceIDs: []string{"ks-123"},
})

// Atualizar
updated, err := client.Agent.Update(ctx, "agent-123", &agentsdk.UpdateAgentRequest{
    Temperature: agentsdk.Float64Ptr(0.7),
    Status:      agentsdk.AgentStatusPtr(agentsdk.AgentStatusInactive),
})

// Deletar
err := client.Agent.Delete(ctx, "agent-123")

// Executar (síncrono)
result, err := client.Agent.Execute(ctx, "agent-123", &agentsdk.ExecuteAgentRequest{
    Input:     "Como faço para resetar minha senha?",
    SessionID: "sessao-001",
    Variables: map[string]any{"user": "João"},
})
fmt.Println(result.Output)
fmt.Printf("Tokens: %d\n", result.Usage.TotalTokens)

// Executar (streaming SSE)
err := client.Agent.ExecuteStream(ctx, "agent-123", &agentsdk.ExecuteAgentRequest{
    Input: "Explique o processo de onboarding.",
}, func(event agentsdk.StreamEvent) error {
    switch event.Type {
    case "delta":
        fmt.Print(event.Delta)  // fragmento de texto
    case "done":
        fmt.Printf("\n[Concluído — %d tokens]\n", event.Usage.TotalTokens)
    case "error":
        return fmt.Errorf("erro: %s", event.Error)
    }
    return nil
})
```

---

## Configuração do Client

```go
client := agentsdk.New(
    "client-id",
    "client-secret",

    // URL base da API
    agentsdk.WithBaseURL("https://api.exemplo.io/v1"),

    // Timeout customizado (padrão: 30s)
    agentsdk.WithTimeout(15 * time.Second),

    // Injetar http.Client customizado (proxies, TLS etc.)
    agentsdk.WithHTTPClient(&http.Client{Transport: myTransport}),

    // Pular autenticação automática usando token existente
    agentsdk.WithToken("Bearer tok_existente"),
)
```

---

## Tratamento de Erros

```go
source, err := client.KnowledgeSource.Get(ctx, "ks-inexistente")
if err != nil {
    switch {
    case agentsdk.IsNotFound(err):
        fmt.Println("Recurso não encontrado")
    case agentsdk.IsUnauthorized(err):
        fmt.Println("Token inválido ou expirado")
    case agentsdk.IsForbidden(err):
        fmt.Println("Sem permissão para este recurso")
    case agentsdk.IsRateLimit(err):
        fmt.Println("Limite de requisições atingido")
    default:
        fmt.Printf("Erro inesperado: %v\n", err)
    }
}

// Inspecionar detalhes do erro
if apiErr, ok := err.(*agentsdk.Error); ok {
    fmt.Printf("HTTP %d | Código: %s | Mensagem: %s\n",
        apiErr.HTTPStatus, apiErr.Code, apiErr.Message)
}
```

---

## Funções Helper

```go
// Ponteiros para tipos primitivos (necessários nos campos *T dos requests de Update)
agentsdk.StringPtr("valor")        // *string
agentsdk.IntPtr(1024)              // *int
agentsdk.Float64Ptr(0.5)          // *float64
agentsdk.BoolPtr(true)            // *bool
agentsdk.AgentStatusPtr(agentsdk.AgentStatusActive)
agentsdk.KnowledgeSourceStatusPtr(agentsdk.KnowledgeSourceStatusActive)
agentsdk.KnowledgeObjectStatusPtr(agentsdk.KnowledgeObjectStatusIndexing)
agentsdk.AgentModelPtr(agentsdk.AgentModelGPT4o)
```

---

## Executando os Testes

```bash
cd agentsdk
go test ./... -v -race
```

---

## Estrutura do Projeto

```
agentsdk/
├── client.go              # Client principal, options, request/response helpers
├── auth.go                # AuthService: GenerateToken, RefreshToken, RevokeToken
├── knowledge_source.go    # KnowledgeSourceService: List, Get, Create, Update, Delete
├── knowledge_object.go    # KnowledgeObjectService: List, Get, Create, Update, Delete
├── agent.go               # AgentService: List, Get, Create, Update, Delete, Execute, ExecuteStream
├── models.go              # Todos os tipos, structs e constantes
├── errors.go              # Error type, helpers IsNotFound, IsUnauthorized etc.
├── helpers.go             # Funções utilitárias (StringPtr, IntPtr etc.)
├── auth_test.go           # Testes de autenticação
├── knowledge_source_test.go
├── agent_test.go          # Testes de Knowledge Object e Agentes
├── examples/
│   └── main.go            # Exemplo completo de uso
└── go.mod
```

---

## Licença

MIT# toolkit-stackspot-ai
