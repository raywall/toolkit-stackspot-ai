package types

import "time"

// ─────────────────────────────────────────────────────────────────────────────
// Knowledge Object models
// ─────────────────────────────────────────────────────────────────────────────

// KnowledgeObjectStatus representa o status do objeto.
type KnowledgeObjectStatus string

const (
	KnowledgeObjectStatusActive   KnowledgeObjectStatus = "active"
	KnowledgeObjectStatusInactive KnowledgeObjectStatus = "inactive"
	KnowledgeObjectStatusIndexing KnowledgeObjectStatus = "indexing"
	KnowledgeObjectStatusError    KnowledgeObjectStatus = "error"
)

// KnowledgeObject representa um objeto de conhecimento dentro de uma Knowledge Source.
type KnowledgeObject struct {
	ID                string                `json:"id"`
	KnowledgeSourceID string                `json:"knowledge_source_id"`
	Title             string                `json:"title"`
	Content           string                `json:"content,omitempty"`
	ContentURL        string                `json:"content_url,omitempty"`
	MimeType          string                `json:"mime_type,omitempty"`
	Status            KnowledgeObjectStatus `json:"status"`
	Tags              []string              `json:"tags,omitempty"`
	Metadata          map[string]any        `json:"metadata,omitempty"`
	EmbeddingModel    string                `json:"embedding_model,omitempty"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
}
