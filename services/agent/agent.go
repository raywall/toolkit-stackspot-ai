package agent

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/raywall/toolkit-stackspot-ai/pkg/clients"
	"github.com/raywall/toolkit-stackspot-ai/pkg/config"
	"github.com/raywall/toolkit-stackspot-ai/pkg/types"
)

// AgentService gerencia as operações de Agentes.
type AgentService struct {
	client *clients.Client
}

// NewAgentService cria uma nova instância do serviço de agentes.
func NewAgentService(c *clients.Client) *AgentService {
	return &AgentService{client: c}
}

// List retorna uma página de Agentes de acordo com os parâmetros fornecidos.
func (s *AgentService) List(ctx context.Context, params *ListAgentsParams) (*types.Page[types.Agent], error) {
	var page types.Page[types.Agent]
	path := config.AgentsBasePath.WithQuery(params.ToMap())

	if err := s.client.DoAuthenticated(ctx, http.MethodGet, path, nil, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// Get recupera um Agente pelo seu ID.
func (s *AgentService) Get(ctx context.Context, id string) (*types.Agent, error) {
	if id == "" {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "id do agente não pode ser vazio"}
	}

	var agent types.Agent
	path := config.AgentsBasePath.Join(id)

	if err := s.client.DoAuthenticated(ctx, http.MethodGet, path, nil, &agent); err != nil {
		return nil, err
	}
	return &agent, nil
}

// Create cria um novo Agente.
func (s *AgentService) Create(ctx context.Context, req *CreateAgentRequest) (*types.Agent, error) {
	if req == nil {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "request não pode ser nil"}
	}
	if req.Name == "" {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "campo 'name' é obrigatório"}
	}
	if req.Model == "" {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "campo 'model' é obrigatório"}
	}

	var agent types.Agent
	if err := s.client.DoAuthenticated(ctx, http.MethodPost, config.AgentsBasePath.String(), req, &agent); err != nil {
		return nil, err
	}
	return &agent, nil
}

// Update atualiza parcialmente um Agente existente.
func (s *AgentService) Update(ctx context.Context, id string, req *UpdateAgentRequest) (*types.Agent, error) {
	if id == "" {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "id do agente não pode ser vazio"}
	}
	if req == nil {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "request não pode ser nil"}
	}

	var agent types.Agent
	path := config.AgentsBasePath.Join(id)

	if err := s.client.DoAuthenticated(ctx, http.MethodPatch, path, req, &agent); err != nil {
		return nil, err
	}
	return &agent, nil
}

// Delete remove permanentemente um Agente.
func (s *AgentService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return &clients.Error{Code: clients.ErrCodeBadRequest, Message: "id do agente não pode ser vazio"}
	}

	path := config.AgentsBasePath.Join(id)
	return s.client.DoAuthenticated(ctx, http.MethodDelete, path, nil, nil)
}

// Execute executa um Agente com a entrada fornecida e retorna a resposta completa.
func (s *AgentService) Execute(ctx context.Context, id string, req *ExecuteAgentRequest) (*ExecuteAgentResponse, error) {
	if id == "" {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "id do agente não pode ser vazio"}
	}
	if req == nil {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "request não pode ser nil"}
	}
	if req.Input == "" {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "campo 'input' é obrigatório"}
	}

	var resp ExecuteAgentResponse
	req.Stream = false // garante modo não-streaming para esta função

	path := config.AgentsBasePath.Join(id, "execute")
	if err := s.client.DoAuthenticated(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ExecuteStream executa um Agente com streaming via Server-Sent Events (SSE).
func (s *AgentService) ExecuteStream(ctx context.Context, id string, req *ExecuteAgentRequest, handler StreamHandler) error {
	if id == "" {
		return &clients.Error{Code: clients.ErrCodeBadRequest, Message: "id do agente não pode ser vazio"}
	}
	if req == nil {
		return &clients.Error{Code: clients.ErrCodeBadRequest, Message: "request não pode ser nil"}
	}
	if req.Input == "" {
		return &clients.Error{Code: clients.ErrCodeBadRequest, Message: "campo 'input' é obrigatório"}
	}
	if handler == nil {
		return &clients.Error{Code: clients.ErrCodeBadRequest, Message: "handler não pode ser nil"}
	}
	if err := s.client.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	req.Stream = true
	path := config.AgentsBasePath.Join(id, "execute")
	httpReq, err := s.client.NewRequest(ctx, http.MethodPost, path, req)
	if err != nil {
		return err
	}
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := s.client.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("erro na requisição de streaming: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var body []byte
		if _, err := fmt.Sscanf(resp.Status, "%d", new(int)); err == nil {
			body = []byte(resp.Status)
		}
		return clients.ParseAPIError(resp.StatusCode, body)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		// Formato SSE: "data: <json>"
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var event StreamEvent
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			continue // ignora linhas malformadas
		}

		if err := handler(event); err != nil {
			return fmt.Errorf("streaming interrompido pelo handler: %w", err)
		}

		if event.Type == "done" || event.Type == "error" {
			break
		}
	}

	return scanner.Err()
}
