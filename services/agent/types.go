package agent

import (
	"strconv"
	"time"

	"github.com/raywall/toolkit-stackspot-ai/pkg/types"
)

// CreateAgentRequest é o payload para criar um novo Agente.
type CreateAgentRequest struct {
	Name               string            `json:"name"`
	Description        string            `json:"description,omitempty"`
	Model              types.AgentModel  `json:"model"`
	SystemPrompt       string            `json:"system_prompt,omitempty"`
	Temperature        float64           `json:"temperature,omitempty"`
	MaxTokens          int               `json:"max_tokens,omitempty"`
	KnowledgeSourceIDs []string          `json:"knowledge_source_ids,omitempty"`
	Tools              []types.AgentTool `json:"tools,omitempty"`
	Tags               []string          `json:"tags,omitempty"`
	Metadata           map[string]any    `json:"metadata,omitempty"`
}

// UpdateAgentRequest é o payload para atualizar um Agente existente.
type UpdateAgentRequest struct {
	Name               *string            `json:"name,omitempty"`
	Description        *string            `json:"description,omitempty"`
	Status             *types.AgentStatus `json:"status,omitempty"`
	Model              *types.AgentModel  `json:"model,omitempty"`
	SystemPrompt       *string            `json:"system_prompt,omitempty"`
	Temperature        *float64           `json:"temperature,omitempty"`
	MaxTokens          *int               `json:"max_tokens,omitempty"`
	KnowledgeSourceIDs []string           `json:"knowledge_source_ids,omitempty"`
	Tools              []types.AgentTool  `json:"tools,omitempty"`
	Tags               []string           `json:"tags,omitempty"`
	Metadata           map[string]any     `json:"metadata,omitempty"`
}

// ListAgentsParams define os filtros para listagem de Agentes.
type ListAgentsParams struct {
	types.Pagination
	Status types.AgentStatus `url:"status,omitempty"`
	Model  types.AgentModel  `url:"model,omitempty"`
	Tag    string            `url:"tag,omitempty"`
	Search string            `url:"search,omitempty"`
}

// ExecuteAgentRequest é o payload para execução de um Agente.
type ExecuteAgentRequest struct {
	Input         string           `json:"input"`
	SessionID     string           `json:"session_id,omitempty"`
	Variables     map[string]any   `json:"variables,omitempty"`
	Stream        bool             `json:"stream,omitempty"`
	OverrideModel types.AgentModel `json:"override_model,omitempty"`
	MaxTokens     int              `json:"max_tokens,omitempty"`
}

// ExecuteAgentResponse é a resposta de uma execução de Agente síncrona.
type ExecuteAgentResponse struct {
	ExecutionID  string                 `json:"execution_id"`
	AgentID      string                 `json:"agent_id"`
	SessionID    string                 `json:"session_id"`
	Output       string                 `json:"output"`
	FinishReason string                 `json:"finish_reason"`
	Usage        *types.TokenUsage      `json:"usage,omitempty"`
	ToolCalls    []types.ToolCallResult `json:"tool_calls,omitempty"`
	Metadata     map[string]any         `json:"metadata,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	DurationMs   int64                  `json:"duration_ms"`
}

// StreamEvent representa um evento recebido durante a execução em streaming.
type StreamEvent struct {
	Type         string                `json:"type"`
	Delta        string                `json:"delta,omitempty"`
	FinishReason string                `json:"finish_reason,omitempty"`
	ToolCall     *types.ToolCallResult `json:"tool_call,omitempty"`
	Usage        *types.TokenUsage     `json:"usage,omitempty"`
	Error        string                `json:"error,omitempty"`
}

// StreamHandler é a função de callback chamada para cada StreamEvent.
type StreamHandler func(event StreamEvent) error

// ToMap convert uma lista de parametros em um mapa de strings
func (p *ListAgentsParams) ToMap() map[string]string {
	var values = make(map[string]string)
	if p.Page > 0 {
		values["page"] = strconv.Itoa(p.Page)
	}
	if p.PageSize > 0 {
		values["page_size"] = strconv.Itoa(p.PageSize)
	}
	if p.Model != "" {
		values["model"] = string(p.Model)
	}
	if p.Tag != "" {
		values["tag"] = p.Tag
	}
	if p.Search != "" {
		values["search"] = p.Search
	}
	return values
}
