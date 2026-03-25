package types

import "time"

// AgentStatus representa o status do agente.
type AgentStatus string

const (
	AgentStatusActive   AgentStatus = "active"
	AgentStatusInactive AgentStatus = "inactive"
	AgentStatusDraft    AgentStatus = "draft"
)

// AgentModel define o modelo de linguagem que o agente utiliza.
type AgentModel string

const (
	AgentModelGPT4o         AgentModel = "gpt-4o"
	AgentModelGPT4          AgentModel = "gpt-4"
	AgentModelClaude3Opus   AgentModel = "claude-3-opus"
	AgentModelClaude3Sonnet AgentModel = "claude-3-sonnet"
	AgentModelGeminiPro     AgentModel = "gemini-pro"
)

// AgentTool representa uma ferramenta disponível ao agente.
type AgentTool struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Config      map[string]any `json:"config,omitempty"`
}

// Agent representa um agente configurado na plataforma.
type Agent struct {
	ID                 string         `json:"id"`
	Name               string         `json:"name"`
	Description        string         `json:"description,omitempty"`
	Status             AgentStatus    `json:"status"`
	Model              AgentModel     `json:"model"`
	SystemPrompt       string         `json:"system_prompt,omitempty"`
	Temperature        float64        `json:"temperature"`
	MaxTokens          int            `json:"max_tokens,omitempty"`
	KnowledgeSourceIDs []string       `json:"knowledge_source_ids,omitempty"`
	Tools              []AgentTool    `json:"tools,omitempty"`
	Tags               []string       `json:"tags,omitempty"`
	Metadata           map[string]any `json:"metadata,omitempty"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

// TokenUsage contém as métricas de uso de tokens de uma execução.
type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ToolCallResult representa o resultado de uma chamada de ferramenta.
type ToolCallResult struct {
	ToolName string         `json:"tool_name"`
	Input    map[string]any `json:"input,omitempty"`
	Output   any            `json:"output,omitempty"`
	Error    string         `json:"error,omitempty"`
}
