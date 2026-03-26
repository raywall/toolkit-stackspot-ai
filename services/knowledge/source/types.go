package source

import (
	"strconv"

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

// ToMap convert uma lista de parametros em um mapa de strings
func (p *ListKnowledgeSourcesParams) ToMap() map[string]string {
	var values = make(map[string]string)
	if p.Page > 0 {
		values["page"] = strconv.Itoa(p.Page)
	}
	if p.PageSize > 0 {
		values["page_size"] = strconv.Itoa(p.PageSize)
	}
	if p.Type != "" {
		values["type"] = string(p.Type)
	}
	if p.Status != "" {
		values["status"] = string(p.Status)
	}
	if p.Tag != "" {
		values["tag"] = p.Tag
	}
	if p.Search != "" {
		values["search"] = p.Search
	}
	return values
}
