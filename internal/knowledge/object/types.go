package object

import (
	"github.com/raywall/toolkit-stackspot-ai/pkg/types"
)

// CreateKnowledgeObjectRequest é o payload para criar um novo Knowledge Object.
type CreateKnowledgeObjectRequest struct {
	Title          string         `json:"title"`
	Content        string         `json:"content,omitempty"`
	ContentURL     string         `json:"content_url,omitempty"`
	MimeType       string         `json:"mime_type,omitempty"`
	Tags           []string       `json:"tags,omitempty"`
	Metadata       map[string]any `json:"metadata,omitempty"`
	EmbeddingModel string         `json:"embedding_model,omitempty"`
}

// UpdateKnowledgeObjectRequest é o payload para atualizar um Knowledge Object.
type UpdateKnowledgeObjectRequest struct {
	Title      *string                      `json:"title,omitempty"`
	Content    *string                      `json:"content,omitempty"`
	ContentURL *string                      `json:"content_url,omitempty"`
	Status     *types.KnowledgeObjectStatus `json:"status,omitempty"`
	Tags       []string                     `json:"tags,omitempty"`
	Metadata   map[string]any               `json:"metadata,omitempty"`
}

// ListKnowledgeObjectsParams define os filtros para listagem de Knowledge Objects.
type ListKnowledgeObjectsParams struct {
	types.Pagination
	Status types.KnowledgeObjectStatus `url:"status,omitempty"`
	Tag    string                      `url:"tag,omitempty"`
	Search string                      `url:"search,omitempty"`
}

// CreateUploadRequest é o payload para iniciar a intenção de upload de um arquivo.
type CreateUploadRequest struct {
	FileName string `json:"file_name"`
	MimeType string `json:"mime_type,omitempty"`
}

// CreateUploadResponse contém o ID gerado para o upload e, opcionalmente, uma URL.
type CreateUploadResponse struct {
	UploadID  string `json:"upload_id"`
	UploadURL string `json:"upload_url,omitempty"` // URL pré-assinada (se a API usar upload direto)
}
