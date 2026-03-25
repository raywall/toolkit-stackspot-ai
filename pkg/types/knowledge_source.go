package types

import "time"

// ─────────────────────────────────────────────────────────────────────────────
// Knowledge Source models
// ─────────────────────────────────────────────────────────────────────────────

// KnowledgeSourceType define o tipo da fonte de conhecimento.
type KnowledgeSourceType string

const (
	KnowledgeSourceTypeDocument  KnowledgeSourceType = "document"
	KnowledgeSourceTypeDatabase  KnowledgeSourceType = "database"
	KnowledgeSourceTypeAPI       KnowledgeSourceType = "api"
	KnowledgeSourceTypeWebScrape KnowledgeSourceType = "web_scrape"
)

// KnowledgeSourceStatus representa o status de processamento da fonte.
type KnowledgeSourceStatus string

const (
	KnowledgeSourceStatusActive     KnowledgeSourceStatus = "active"
	KnowledgeSourceStatusInactive   KnowledgeSourceStatus = "inactive"
	KnowledgeSourceStatusProcessing KnowledgeSourceStatus = "processing"
	KnowledgeSourceStatusError      KnowledgeSourceStatus = "error"
)

// KnowledgeSource representa uma fonte de conhecimento na plataforma.
type KnowledgeSource struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description,omitempty"`
	Type        KnowledgeSourceType   `json:"type"`
	Status      KnowledgeSourceStatus `json:"status"`
	Config      map[string]any        `json:"config,omitempty"`
	Tags        []string              `json:"tags,omitempty"`
	Metadata    map[string]any        `json:"metadata,omitempty"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}
