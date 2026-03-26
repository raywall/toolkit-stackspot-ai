package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/raywall/toolkit-stackspot-ai/pkg/config"
)

const (
	defaultTimeout    = 30 * time.Second
	headerAuthToken   = "Authorization"
	headerContentType = "Content-Type"
	contentTypeJSON   = "application/json"
)

// Client é o motor HTTP principal da SDK.
type Client struct {
	HTTPClient  *http.Client
	AuthBaseURL string
	BaseURL     string
	Realm       string
	Token       string

	// TokenProvider é chamado quando uma requisição autenticada é feita
	// mas o Token atual está vazio. Deve ser configurado pelo orquestrador.
	TokenProvider func(ctx context.Context) error
}

// Option é uma função de configuração usada para personalizar o Client.
type Option func(*Client)

// WithHTTPClient substitui o http.Client padrão.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.HTTPClient = hc }
}

// WithBaseURL define a URL base da API.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) { c.BaseURL = baseURL }
}

// WithTimeout define o timeout das requisições HTTP.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.HTTPClient.Timeout = d }
}

// WithToken injeta diretamente um token de autenticação,
// pulando a etapa de autenticação automática.
func WithToken(token string) Option {
	return func(c *Client) { c.Token = token }
}

// New cria um novo Client HTTP base.
func New(opts ...Option) *Client {
	c := &Client{
		HTTPClient:  &http.Client{Timeout: defaultTimeout},
		AuthBaseURL: config.AuthBaseUrl.String(),
		BaseURL:     config.ApiBaseUrl.String(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// EnsureAuthenticated garante que o client possui um token,
// acionando o TokenProvider caso o token esteja vazio.
func (c *Client) EnsureAuthenticated(ctx context.Context) error {
	if c.Token != "" {
		return nil
	}
	if c.TokenProvider != nil {
		return c.TokenProvider(ctx)
	}
	return &Error{Code: ErrCodeUnauthorized, Message: "token ausente e provider não configurado"}
}

// NewRequest constrói um *http.Request com os headers padrão.
// O parâmetro path pode conter query string (ex: "/resources?page=1"),
// que será preservada corretamente na URL final.
func (c *Client) NewRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	rawPath := path
	rawQuery := ""
	if idx := strings.IndexByte(path, '?'); idx >= 0 {
		rawPath = path[:idx]
		rawQuery = path[idx+1:]
	}

	joined, err := url.JoinPath(c.BaseURL, rawPath)
	if err != nil {
		return nil, fmt.Errorf("url inválida: %w", err)
	}

	finalURL := joined
	if rawQuery != "" {
		finalURL = joined + "?" + rawQuery
	}

	var buf io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("erro ao serializar body: %w", err)
		}
		buf = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, finalURL, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set(headerContentType, contentTypeJSON)
	if c.Token != "" {
		req.Header.Set(headerAuthToken, "Bearer "+c.Token)
	}
	return req, nil
}

// Do executa a requisição HTTP e decodifica a resposta JSON em out.
func (c *Client) Do(req *http.Request, out any) error {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("erro na requisição HTTP: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erro ao ler resposta: %w", err)
	}

	if resp.StatusCode >= 400 {
		return ParseAPIError(resp.StatusCode, bodyBytes)
	}

	if out != nil && len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, out); err != nil {
			return fmt.Errorf("erro ao desserializar resposta: %w", err)
		}
	}
	return nil
}

// DoAuthenticated executa uma requisição autenticada, garantindo token válido.
func (c *Client) DoAuthenticated(ctx context.Context, method, path string, body, out any) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}
	req, err := c.NewRequest(ctx, method, path, body)
	if err != nil {
		return err
	}
	return c.Do(req, out)
}
