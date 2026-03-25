package object

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/raywall/toolkit-stackspot-ai/pkg/clients"
	"github.com/raywall/toolkit-stackspot-ai/pkg/types"
)

// KnowledgeObjectService gerencia as operações de Knowledge Objects
// dentro de uma Knowledge Source.
type KnowledgeObjectService struct {
	client *clients.Client
}

// NewKnowledgeObjectService cria uma nova instância do serviço.
func NewKnowledgeObjectService(c *clients.Client) *KnowledgeObjectService {
	return &KnowledgeObjectService{client: c}
}

func knowledgeObjectsPath(sourceID string) string {
	return fmt.Sprintf("/knowledge-sources/%s/objects", sourceID)
}

func knowledgeObjectPath(sourceID, objectID string) string {
	return fmt.Sprintf("/knowledge-sources/%s/objects/%s", sourceID, objectID)
}

// List retorna uma página de Knowledge Objects de uma Knowledge Source específica.
func (s *KnowledgeObjectService) List(ctx context.Context, sourceID string, params *ListKnowledgeObjectsParams) (*types.Page[types.KnowledgeObject], error) {
	if sourceID == "" {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "sourceID não pode ser vazio"}
	}

	path := knowledgeObjectsPath(sourceID)
	if params != nil {
		path += "?" + encodeKnowledgeObjectParams(params)
	}

	var page types.Page[types.KnowledgeObject]
	if err := s.client.DoAuthenticated(ctx, http.MethodGet, path, nil, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// Get recupera um Knowledge Object específico pelo ID.
func (s *KnowledgeObjectService) Get(ctx context.Context, sourceID, objectID string) (*types.KnowledgeObject, error) {
	if sourceID == "" {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "sourceID não pode ser vazio"}
	}
	if objectID == "" {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "objectID não pode ser vazio"}
	}

	var ko types.KnowledgeObject
	if err := s.client.DoAuthenticated(ctx, http.MethodGet, knowledgeObjectPath(sourceID, objectID), nil, &ko); err != nil {
		return nil, err
	}
	return &ko, nil
}

// Create cria um novo Knowledge Object dentro de uma Knowledge Source.
func (s *KnowledgeObjectService) Create(ctx context.Context, sourceID string, req *CreateKnowledgeObjectRequest) (*types.KnowledgeObject, error) {
	if sourceID == "" {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "sourceID não pode ser vazio"}
	}
	if req == nil {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "request não pode ser nil"}
	}
	if req.Title == "" {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "campo 'title' é obrigatório"}
	}
	if req.Content == "" && req.ContentURL == "" {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "pelo menos um de 'content' ou 'content_url' deve ser fornecido"}
	}

	var ko types.KnowledgeObject
	if err := s.client.DoAuthenticated(ctx, http.MethodPost, knowledgeObjectsPath(sourceID), req, &ko); err != nil {
		return nil, err
	}
	return &ko, nil
}

// Update atualiza parcialmente um Knowledge Object existente.
func (s *KnowledgeObjectService) Update(ctx context.Context, sourceID, objectID string, req *UpdateKnowledgeObjectRequest) (*types.KnowledgeObject, error) {
	if sourceID == "" {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "sourceID não pode ser vazio"}
	}
	if objectID == "" {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "objectID não pode ser vazio"}
	}
	if req == nil {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "request não pode ser nil"}
	}

	var ko types.KnowledgeObject
	if err := s.client.DoAuthenticated(ctx, http.MethodPatch, knowledgeObjectPath(sourceID, objectID), req, &ko); err != nil {
		return nil, err
	}
	return &ko, nil
}

// Delete remove um Knowledge Object permanentemente.
func (s *KnowledgeObjectService) Delete(ctx context.Context, sourceID, objectID string) error {
	if sourceID == "" {
		return &clients.Error{Code: clients.ErrCodeBadRequest, Message: "sourceID não pode ser vazio"}
	}
	if objectID == "" {
		return &clients.Error{Code: clients.ErrCodeBadRequest, Message: "objectID não pode ser vazio"}
	}

	return s.client.DoAuthenticated(ctx, http.MethodDelete, knowledgeObjectPath(sourceID, objectID), nil, nil)
}

// encodeKnowledgeObjectParams converte ListKnowledgeObjectsParams em query string.
func encodeKnowledgeObjectParams(p *ListKnowledgeObjectsParams) string {
	v := url.Values{}
	if p.Page > 0 {
		v.Set("page", strconv.Itoa(p.Page))
	}
	if p.PageSize > 0 {
		v.Set("page_size", strconv.Itoa(p.PageSize))
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

// GenerateUpload cria uma intenção de upload para um arquivo local, retornando o Upload ID.
func (s *KnowledgeObjectService) GenerateUpload(ctx context.Context, sourceID string, req *CreateUploadRequest) (*CreateUploadResponse, error) {
	if sourceID == "" {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "sourceID não pode ser vazio"}
	}
	if req == nil || req.FileName == "" {
		return nil, &clients.Error{Code: clients.ErrCodeBadRequest, Message: "request nulo ou file_name vazio"}
	}

	path := fmt.Sprintf("/knowledge-sources/%s/objects/upload", sourceID)
	var resp CreateUploadResponse
	if err := s.client.DoAuthenticated(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UploadFile realiza o upload do arquivo físico usando multipart/form-data a partir
// do diretório (outFileDir) e nome do arquivo (outFileName).
func (s *KnowledgeObjectService) UploadFile(ctx context.Context, sourceID, uploadID, outFileDir, outFileName string) error {
	if sourceID == "" || uploadID == "" {
		return &clients.Error{Code: clients.ErrCodeBadRequest, Message: "sourceID e uploadID são obrigatórios"}
	}

	// 1. Abre o arquivo físico
	filePath := filepath.Join(outFileDir, outFileName)
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo para upload: %w", err)
	}
	defer file.Close()

	// 2. Prepara o buffer multipart
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", outFileName)
	if err != nil {
		return fmt.Errorf("erro ao criar form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("erro ao copiar conteúdo do arquivo: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("erro ao fechar writer do multipart: %w", err)
	}

	// 3. Monta a requisição reaproveitando o client base para pegar BaseURL e Token
	path := fmt.Sprintf("/knowledge-sources/%s/objects/%s/file", sourceID, uploadID)

	if err := s.client.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	// Passamos body nulo no NewRequest para evitar que ele serialize como JSON
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return err
	}

	// 4. Sobrescrevemos os headers e o body com os dados do multipart
	req.Body = io.NopCloser(body)
	req.ContentLength = int64(body.Len())
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 5. Executa a requisição
	return s.client.Do(req, nil)
}

// DeleteAll remove todos os Knowledge Objects atrelados a uma Knowledge Source.
func (s *KnowledgeObjectService) DeleteAll(ctx context.Context, sourceID string) error {
	if sourceID == "" {
		return &clients.Error{Code: clients.ErrCodeBadRequest, Message: "sourceID não pode ser vazio"}
	}

	// Assumindo que o endpoint DELETE base da listagem remova todos os objetos
	path := knowledgeObjectsPath(sourceID)
	return s.client.DoAuthenticated(ctx, http.MethodDelete, path, nil, nil)
}
