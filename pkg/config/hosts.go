package config

import (
	"strings"
)

type (
	BaseUrl  string
	BasePath string
)

const (
	// urls base utilizadas na stackspot
	AuthBaseUrl BaseUrl = "https://idm.stackspot.com"
	ApiBaseUrl  BaseUrl = "https://data-integration-api.stackspot.com"

	// rotas das funcionalidades da stackspot
	AuthBasePath               BasePath = "/oidc/oauth/token"
	AuthRefreshBasePath        BasePath = "/oidc/oauth/token/refresh"
	AuthRevokeBasePath         BasePath = "/oidc/oauth/token/revoke"
	AgentsBasePath             BasePath = "/v1/agents"
	KnowledgeSourcesBasePathV1 BasePath = "/v1/knowledge-sources"
	KnowledgeSourcesBasePathV2 BasePath = "/v2/knowledge-sources"
)

func (b BaseUrl) String() string {
	return string(b)
}

func (b BasePath) String() string {
	return string(b)
}

func (b BasePath) Join(routes ...string) string {
	if len(routes) == 0 {
		return b.String()
	}
	values := []string{b.String()}
	for _, route := range routes {
		values = append(values, route)
	}
	return strings.Join(values, "/")
}

func (b BasePath) WithQuery(query map[string]string) string {
	if len(query) == 0 {
		return b.String()
	}
	values := []string{b.String()}
	for k, v := range query {
		values = append(values, k+"="+v)
	}
	data := strings.Join(values, "&")
	if data != "" {
		data = "?" + data
	}
	return data
}
