package source

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/raywall/toolkit-stackspot-ai/pkg/clients"
	"github.com/raywall/toolkit-stackspot-ai/pkg/types"
)

const knowledgeSourcesBasePath = "/knowledge-sources"

// KnowledgeSourceService gerencia as operações de Knowledge Sources.
type KnowledgeSourceService struct {
	client *clients.Client
}

// NewKnowledgeSourceService cria uma nova instância do serviço de Knowledge Sources.
func NewKnowledgeSourceService(c *clients.Client) *KnowledgeSourceService {
	return &KnowledgeSourceService{client: c}
}

// List retorna uma página de Knowledge Sources de acordo com os parâmetros fornecidos.
func (s *KnowledgeSourceService) List(ctx context.Context, params *ListKnowledgeSourcesParams) (*types.Page[types.KnowledgeSource], error) {
	path := knowledgeSourcesBasePath
	if params != nil {
		path += "?" + encodeKnowledgeSourceParams(params)
	}

	var page types.Page[types.KnowledgeSource]
	if err := s.client.DoAuthenticated(ctx, http.MethodGet, path, nil, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// Get recupera uma Knowledge Source pelo seu ID.
func (s *KnowledgeSourceService) Get(ctx context.Context, id string) (*types.KnowledgeSource, error) {
	if id == "" {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "id da knowledge source não pode ser vazio"}
	}

	path := fmt.Sprintf("%s/%s", knowledgeSourcesBasePath, id)
	var ks types.KnowledgeSource
	if err := s.client.DoAuthenticated(ctx, http.MethodGet, path, nil, &ks); err != nil {
		return nil, err
	}
	return &ks, nil
}

// Create cria uma nova Knowledge Source.
func (s *KnowledgeSourceService) Create(ctx context.Context, req *CreateKnowledgeSourceRequest) (*types.KnowledgeSource, error) {
	if req == nil {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "request não pode ser nil"}
	}
	if req.Name == "" {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "campo 'name' é obrigatório"}
	}
	if req.Type == "" {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "campo 'type' é obrigatório"}
	}

	var ks types.KnowledgeSource
	if err := s.client.DoAuthenticated(ctx, http.MethodPost, knowledgeSourcesBasePath, req, &ks); err != nil {
		return nil, err
	}
	return &ks, nil
}

// Update atualiza parcialmente uma Knowledge Source existente.
// Apenas os campos não-nulos do UpdateKnowledgeSourceRequest serão alterados.
func (s *KnowledgeSourceService) Update(ctx context.Context, id string, req *UpdateKnowledgeSourceRequest) (*types.KnowledgeSource, error) {
	if id == "" {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "id da knowledge source não pode ser vazio"}
	}
	if req == nil {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "request não pode ser nil"}
	}

	path := fmt.Sprintf("%s/%s", knowledgeSourcesBasePath, id)
	var ks types.KnowledgeSource
	if err := s.client.DoAuthenticated(ctx, http.MethodPatch, path, req, &ks); err != nil {
		return nil, err
	}
	return &ks, nil
}

// Delete remove permanentemente uma Knowledge Source e todos os seus objetos.
func (s *KnowledgeSourceService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return &clients.Error{Code: clients.ErrCodeBadRequest, Message: "id da knowledge source não pode ser vazio"}
	}

	path := fmt.Sprintf("%s/%s", knowledgeSourcesBasePath, id)
	return s.client.DoAuthenticated(ctx, http.MethodDelete, path, nil, nil)
}

// encodeKnowledgeSourceParams converte ListKnowledgeSourcesParams em query string.
func encodeKnowledgeSourceParams(p *ListKnowledgeSourcesParams) string {
	v := url.Values{}
	if p.Page > 0 {
		v.Set("page", strconv.Itoa(p.Page))
	}
	if p.PageSize > 0 {
		v.Set("page_size", strconv.Itoa(p.PageSize))
	}
	if p.Type != "" {
		v.Set("type", string(p.Type))
	}
	if p.Status != "" {
		v.Set("status", string(p.Status))
	}
	if p.Tag != "" {
		v.Set("tag", p.Tag)
	}
	if p.Search != "" {
		v.Set("search", p.Search)
	}
	return v.Encode()
}
