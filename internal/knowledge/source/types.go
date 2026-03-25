package source

import (
	"github.com/raywall/toolkit-stackspot-ai/pkg/types"
)

// CreateKnowledgeSourceRequest é o payload para criar uma nova Knowledge Source.
type CreateKnowledgeSourceRequest struct {
	Name        string                    `json:"name"`
	Description string                    `json:"description,omitempty"`
	Type        types.KnowledgeSourceType `json:"type"`
	Config      map[string]any            `json:"config,omitempty"`
	Tags        []string                  `json:"tags,omitempty"`
	Metadata    map[string]any            `json:"metadata,omitempty"`
}

// UpdateKnowledgeSourceRequest é o payload para atualizar uma Knowledge Source existente.
// Campos omitidos (zero-value) não serão alterados.
type UpdateKnowledgeSourceRequest struct {
	Name        *string                      `json:"name,omitempty"`
	Description *string                      `json:"description,omitempty"`
	Status      *types.KnowledgeSourceStatus `json:"status,omitempty"`
	Config      map[string]any               `json:"config,omitempty"`
	Tags        []string                     `json:"tags,omitempty"`
	Metadata    map[string]any               `json:"metadata,omitempty"`
}

// ListKnowledgeSourcesParams define os filtros para listagem de Knowledge Sources.
type ListKnowledgeSourcesParams struct {
	types.Pagination
	Type   types.KnowledgeSourceType   `url:"type,omitempty"`
	Status types.KnowledgeSourceStatus `url:"status,omitempty"`
	Tag    string                      `url:"tag,omitempty"`
	Search string                      `url:"search,omitempty"`
}
