package object_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/raywall/toolkit-stackspot-ai/internal/knowledge/object"
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

func makeKO(id, sourceID, title string) types.KnowledgeObject {
	return types.KnowledgeObject{
		ID:                id,
		KnowledgeSourceID: sourceID,
		Title:             title,
		Content:           "Conteúdo de exemplo",
		Status:            types.KnowledgeObjectStatusActive,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}

func TestKnowledgeObjectService_List(t *testing.T) {
	ctx := context.Background()
	page := types.Page[types.KnowledgeObject]{
		Items:      []types.KnowledgeObject{makeKO("ko-1", "ks-1", "Doc A"), makeKO("ko-2", "ks-1", "Doc B")},
		TotalItems: 2,
		Page:       1,
		PageSize:   20,
	}

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/knowledge-sources/ks-1/objects" && r.Method == http.MethodGet {
			writeJSON(w, http.StatusOK, page)
			return
		}
		http.NotFound(w, r)
	}))

	svc := object.NewKnowledgeObjectService(client)
	result, err := svc.List(ctx, "ks-1", nil)
	if err != nil {
		t.Fatalf("List falhou: %v", err)
	}
	if len(result.Items) != 2 {
		t.Errorf("esperado 2 itens, obteve %d", len(result.Items))
	}
}

func TestKnowledgeObjectService_Create(t *testing.T) {
	ctx := context.Background()
	created := makeKO("ko-new", "ks-1", "Novo Objeto")

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/knowledge-sources/ks-1/objects" && r.Method == http.MethodPost {
			writeJSON(w, http.StatusCreated, created)
			return
		}
		http.NotFound(w, r)
	}))

	svc := object.NewKnowledgeObjectService(client)
	result, err := svc.Create(ctx, "ks-1", &object.CreateKnowledgeObjectRequest{
		Title:   "Novo Objeto",
		Content: "Conteúdo do novo objeto",
	})
	if err != nil {
		t.Fatalf("Create falhou: %v", err)
	}
	if result.ID != "ko-new" {
		t.Errorf("ID esperado ko-new, obteve %s", result.ID)
	}
}

func TestKnowledgeObjectService_Create_Validation(t *testing.T) {
	client := clients.New(clients.WithBaseURL("http://localhost"))
	svc := object.NewKnowledgeObjectService(client)
	ctx := context.Background()

	// Sem title
	_, err := svc.Create(ctx, "ks-1", &object.CreateKnowledgeObjectRequest{Content: "teste"})
	if err == nil {
		t.Error("esperado erro para title vazio")
	}

	// Sem content e sem content_url
	_, err = svc.Create(ctx, "ks-1", &object.CreateKnowledgeObjectRequest{Title: "Teste"})
	if err == nil {
		t.Error("esperado erro quando content e content_url são ambos vazios")
	}
}

func TestKnowledgeObjectService_Delete(t *testing.T) {
	ctx := context.Background()
	deleted := false

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/knowledge-sources/ks-1/objects/ko-1" && r.Method == http.MethodDelete {
			deleted = true
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.NotFound(w, r)
	}))

	svc := object.NewKnowledgeObjectService(client)
	if err := svc.Delete(ctx, "ks-1", "ko-1"); err != nil {
		t.Fatalf("Delete falhou: %v", err)
	}
	if !deleted {
		t.Error("endpoint DELETE não foi chamado")
	}
}
