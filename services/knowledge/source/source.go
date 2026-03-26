package source

import (
	"context"
	"net/http"

	"github.com/raywall/toolkit-stackspot-ai/pkg/clients"
	"github.com/raywall/toolkit-stackspot-ai/pkg/config"
	"github.com/raywall/toolkit-stackspot-ai/pkg/types"
)

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
	var page types.Page[types.KnowledgeSource]
	path := config.KnowledgeSourcesBasePathV1.WithQuery(params.ToMap())

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

	var ks types.KnowledgeSource
	path := config.KnowledgeSourcesBasePathV2.Join(id)

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
	if err := s.client.DoAuthenticated(ctx, http.MethodPost, config.KnowledgeSourcesBasePathV1.String(), req, &ks); err != nil {
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

	var ks types.KnowledgeSource
	path := config.KnowledgeSourcesBasePathV2.Join(id)

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

	path := config.KnowledgeSourcesBasePathV2.Join(id)
	return s.client.DoAuthenticated(ctx, http.MethodDelete, path, nil, nil)
}
