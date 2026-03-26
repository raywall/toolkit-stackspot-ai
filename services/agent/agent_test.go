package agent_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/raywall/toolkit-stackspot-ai/pkg/clients"
	"github.com/raywall/toolkit-stackspot-ai/pkg/config"
	"github.com/raywall/toolkit-stackspot-ai/pkg/types"
	"github.com/raywall/toolkit-stackspot-ai/services/agent"
)

func newTestClient(t *testing.T, handler http.Handler) *clients.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return clients.New(
		clients.WithBaseURL(srv.URL),
		clients.WithToken("tok"), // Bypass na necessidade de auth
	)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func makeAgent(id, name string) types.Agent {
	return types.Agent{
		ID:          id,
		Name:        name,
		Status:      types.AgentStatusActive,
		Model:       types.AgentModelGPT4o,
		Temperature: 0.5,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func TestAgentService_List(t *testing.T) {
	ctx := context.Background()
	page := types.Page[types.Agent]{
		Items:      []types.Agent{makeAgent("ag-1", "Agente A"), makeAgent("ag-2", "Agente B")},
		TotalItems: 2,
		Page:       1,
		PageSize:   10,
	}

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == config.AgentsBasePath.String() && r.Method == http.MethodGet {
			writeJSON(w, http.StatusOK, page)
			return
		}
		http.NotFound(w, r)
	}))

	svc := agent.NewAgentService(client)
	result, err := svc.List(ctx, &agent.ListAgentsParams{
		Pagination: types.Pagination{Page: 1, PageSize: 10},
		Status:     types.AgentStatusActive,
	})
	if err != nil {
		t.Fatalf("List falhou: %v", err)
	}
	if len(result.Items) != 2 {
		t.Errorf("esperado 2 agentes, obteve %d", len(result.Items))
	}
}

func TestAgentService_Create(t *testing.T) {
	ctx := context.Background()
	created := makeAgent("ag-new", "Meu Agente")

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == config.AgentsBasePath.String() && r.Method == http.MethodPost {
			writeJSON(w, http.StatusCreated, created)
			return
		}
		http.NotFound(w, r)
	}))

	svc := agent.NewAgentService(client)
	result, err := svc.Create(ctx, &agent.CreateAgentRequest{
		Name:         "Meu Agente",
		Model:        types.AgentModelGPT4o,
		SystemPrompt: "Você é um assistente prestativo.",
		Temperature:  0.5,
	})
	if err != nil {
		t.Fatalf("Create falhou: %v", err)
	}
	if result.ID != "ag-new" {
		t.Errorf("ID esperado ag-new, obteve %s", result.ID)
	}
}

func TestAgentService_Update(t *testing.T) {
	ctx := context.Background()
	updated := makeAgent("ag-1", "Agente Atualizado")

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == config.AgentsBasePath.Join("ag-1") && r.Method == http.MethodPatch {
			writeJSON(w, http.StatusOK, updated)
			return
		}
		http.NotFound(w, r)
	}))

	svc := agent.NewAgentService(client)
	result, err := svc.Update(ctx, "ag-1", &agent.UpdateAgentRequest{
		Name: types.StringPtr("Agente Atualizado"),
	})
	if err != nil {
		t.Fatalf("Update falhou: %v", err)
	}
	if result.Name != "Agente Atualizado" {
		t.Errorf("esperado 'Agente Atualizado', obteve %s", result.Name)
	}
}

func TestAgentService_Delete(t *testing.T) {
	ctx := context.Background()
	deleted := false

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == config.AgentsBasePath.Join("ag-1") && r.Method == http.MethodDelete {
			deleted = true
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.NotFound(w, r)
	}))

	svc := agent.NewAgentService(client)
	if err := svc.Delete(ctx, "ag-1"); err != nil {
		t.Fatalf("Delete falhou: %v", err)
	}
	if !deleted {
		t.Error("endpoint DELETE não foi chamado")
	}
}

func TestAgentService_Execute(t *testing.T) {
	ctx := context.Background()
	execResp := agent.ExecuteAgentResponse{
		ExecutionID:  "exec-001",
		AgentID:      "ag-1",
		SessionID:    "sess-abc",
		Output:       "Olá! Posso ajudá-lo com suporte técnico.",
		FinishReason: "stop",
		Usage: &types.TokenUsage{
			PromptTokens:     50,
			CompletionTokens: 30,
			TotalTokens:      80,
		},
		DurationMs: 1200,
	}

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == config.AgentsBasePath.Join("ag-1", "execute") && r.Method == http.MethodPost {
			writeJSON(w, http.StatusOK, execResp)
			return
		}
		http.NotFound(w, r)
	}))

	svc := agent.NewAgentService(client)
	result, err := svc.Execute(ctx, "ag-1", &agent.ExecuteAgentRequest{
		Input:     "Como faço para resetar minha senha?",
		SessionID: "sess-abc",
	})
	if err != nil {
		t.Fatalf("Execute falhou: %v", err)
	}
	if result.ExecutionID != "exec-001" {
		t.Errorf("ExecutionID esperado exec-001, obteve %s", result.ExecutionID)
	}
	if result.Output == "" {
		t.Error("Output não deve ser vazio")
	}
	if result.Usage == nil {
		t.Error("Usage não deve ser nil")
	}
	if result.Usage.TotalTokens != 80 {
		t.Errorf("TotalTokens esperado 80, obteve %d", result.Usage.TotalTokens)
	}
}

func TestAgentService_ExecuteStream(t *testing.T) {
	ctx := context.Background()

	sseBody := "data: {\"type\":\"delta\",\"delta\":\"Olá \"}\n\n" +
		"data: {\"type\":\"delta\",\"delta\":\"mundo!\"}\n\n" +
		"data: {\"type\":\"done\",\"finish_reason\":\"stop\",\"usage\":{\"prompt_tokens\":10,\"completion_tokens\":5,\"total_tokens\":15}}\n\n" +
		"data: [DONE]\n\n"

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == config.AgentsBasePath.Join("ag-1", "execute") && r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, sseBody)
			return
		}
		http.NotFound(w, r)
	}))

	var collected strings.Builder
	eventCount := 0

	svc := agent.NewAgentService(client)
	err := svc.ExecuteStream(ctx, "ag-1", &agent.ExecuteAgentRequest{
		Input: "Diga olá",
	}, func(event agent.StreamEvent) error {
		eventCount++
		if event.Type == "delta" {
			collected.WriteString(event.Delta)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("ExecuteStream falhou: %v", err)
	}
	if collected.String() != "Olá mundo!" {
		t.Errorf("output esperado 'Olá mundo!', obteve %q", collected.String())
	}
	if eventCount < 2 {
		t.Errorf("esperado ao menos 2 eventos, obteve %d", eventCount)
	}
}
