package types

// Credentials contém as credenciais do cliente para autenticação.
type Credentials struct {
	// ClientID é o identificador do cliente (client_id OAuth2).
	ClientID string `json:"client_id"`
	// ClientSecret é o segredo do cliente (client_secret OAuth2).
	ClientSecret string `json:"client_secret"`
	// Scope define os escopos solicitados (opcional).
	Scope string `json:"scope,omitempty"`
	// GrantType define o fluxo OAuth2. Padrão: "client_credentials".
	GrantType string `json:"grant_type,omitempty"`
}
