package types

// Pagination encapsula os parâmetros de paginação para listagens.
type Pagination struct {
	// Page é o número da página (base 1).
	Page int `url:"page,omitempty"`
	// PageSize é a quantidade de itens por página.
	PageSize int `url:"page_size,omitempty"`
}

// Page é a resposta paginada genérica.
type Page[T any] struct {
	Items      []T  `json:"items"`
	TotalItems int  `json:"total_items"`
	TotalPages int  `json:"total_pages"`
	Page       int  `json:"page"`
	PageSize   int  `json:"page_size"`
	HasNext    bool `json:"has_next"`
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers (Ponteiros para Tipos Primitivos e de Domínio)
// ─────────────────────────────────────────────────────────────────────────────

// StringPtr retorna um ponteiro para a string fornecida.
func StringPtr(s string) *string { return &s }

// IntPtr retorna um ponteiro para o int fornecido.
func IntPtr(i int) *int { return &i }

// Float64Ptr retorna um ponteiro para o float64 fornecido.
func Float64Ptr(f float64) *float64 { return &f }

// BoolPtr retorna um ponteiro para o bool fornecido.
func BoolPtr(b bool) *bool { return &b }

// AgentStatusPtr retorna um ponteiro para um AgentStatus.
func AgentStatusPtr(s AgentStatus) *AgentStatus { return &s }

// KnowledgeSourceStatusPtr retorna um ponteiro para um KnowledgeSourceStatus.
func KnowledgeSourceStatusPtr(s KnowledgeSourceStatus) *KnowledgeSourceStatus { return &s }

// KnowledgeObjectStatusPtr retorna um ponteiro para um KnowledgeObjectStatus.
func KnowledgeObjectStatusPtr(s KnowledgeObjectStatus) *KnowledgeObjectStatus { return &s }

// AgentModelPtr retorna um ponteiro para um AgentModel.
func AgentModelPtr(m AgentModel) *AgentModel { return &m }
