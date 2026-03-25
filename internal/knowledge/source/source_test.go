package source_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/raywall/toolkit-stackspot-ai/internal/knowledge/source"
	"github.com/raywall/toolkit-stackspot-ai/pkg/clients"
	"github.com/raywall/toolkit-stackspot-ai/pkg/types"
)

func newTestClient(t *testing.T, handler http.Handler) *clients.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return clients.New(
		clients.WithBaseURL(srv.URL),
		clients.WithToken("tok"),
	)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func makeKS(id, name string) types.KnowledgeSource {
	return types.KnowledgeSource{
		ID:        id,
		Name:      name,
		Type:      types.KnowledgeSourceTypeDocument,
		Status:    types.KnowledgeSourceStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestKnowledgeSourceService_List(t *testing.T) {
	ctx := context.Background()
	page := types.Page[types.KnowledgeSource]{
		Items:      []types.KnowledgeSource{makeKS("ks-1", "Fonte A"), makeKS("ks-2", "Fonte B")},
		TotalItems: 2,
		Page:       1,
		PageSize:   20,
	}

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/knowledge-sources" && r.Method == http.MethodGet {
			writeJSON(w, http.StatusOK, page)
			return
		}
		http.NotFound(w, r)
	}))

	svc := source.NewKnowledgeSourceService(client)
	result, err := svc.List(ctx, &source.ListKnowledgeSourcesParams{
		Pagination: types.Pagination{Page: 1, PageSize: 20},
	})
	if err != nil {
		t.Fatalf("List falhou: %v", err)
	}
	if len(result.Items) != 2 {
		t.Errorf("esperado 2 itens, obteve %d", len(result.Items))
	}
}

func TestKnowledgeSourceService_Create(t *testing.T) {
	ctx := context.Background()
	created := makeKS("ks-new", "Nova Fonte")

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/knowledge-sources" && r.Method == http.MethodPost {
			var req source.CreateKnowledgeSourceRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			if req.Name == "" {
				writeJSON(w, http.StatusBadRequest, map[string]string{"message": "name required"})
				return
			}
			writeJSON(w, http.StatusCreated, created)
			return
		}
		http.NotFound(w, r)
	}))

	svc := source.NewKnowledgeSourceService(client)
	result, err := svc.Create(ctx, &source.CreateKnowledgeSourceRequest{
		Name: "Nova Fonte",
		Type: types.KnowledgeSourceTypeDocument,
	})
	if err != nil {
		t.Fatalf("Create falhou: %v", err)
	}
	if result.ID != "ks-new" {
		t.Errorf("ID esperado ks-new, obteve %s", result.ID)
	}
}

func TestKnowledgeSourceService_Update(t *testing.T) {
	ctx := context.Background()
	updated := makeKS("ks-1", "Nome Atualizado")

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/knowledge-sources/ks-1" && r.Method == http.MethodPatch {
			writeJSON(w, http.StatusOK, updated)
			return
		}
		http.NotFound(w, r)
	}))

	svc := source.NewKnowledgeSourceService(client)
	result, err := svc.Update(ctx, "ks-1", &source.UpdateKnowledgeSourceRequest{
		Name: types.StringPtr("Nome Atualizado"),
	})
	if err != nil {
		t.Fatalf("Update falhou: %v", err)
	}
	if result.Name != "Nome Atualizado" {
		t.Errorf("Name esperado 'Nome Atualizado', obteve %s", result.Name)
	}
}

func TestKnowledgeSourceService_NotFound(t *testing.T) {
	ctx := context.Background()
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/knowledge-sources/") {
			writeJSON(w, http.StatusNotFound, map[string]string{
				"code":    "NOT_FOUND",
				"message": "knowledge source não encontrada",
			})
			return
		}
		http.NotFound(w, r)
	}))

	svc := source.NewKnowledgeSourceService(client)
	_, err := svc.Get(ctx, "ks-inexistente")
	if !clients.IsNotFound(err) {
		t.Errorf("esperado erro NotFound, obteve: %v", err)
	}
}
