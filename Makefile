.PHONY: run test coverage

# Captura o segundo argumento da linha de comando (ex: auth, agent, upload)
SAMPLE := $(word 2, $(MAKECMDGOALS))

run:
	@AUTH_BASE_URL=https://idm.stackspot.com/v1 \
	 API_BASE_URL=https://data-integration-api.stackspot.com/v1 \
	 API_CLIENT_ID=3a0541eb-993b-4134-9c93-68d2a9231f55 \
	 API_CLIENT_SECRET=750c9325-d96b-4bff-81bb-79312abe9856 \
	 AUTH_REALM=stackspot \
	 go run samples/$(SAMPLE)/main.go

test:
	@go test -v ./...

coverage:
	@go test -coverprofile=coverage.out ./...; \
	 go tool cover -html=coverage.out -o coverage.html;

# Evita que o Make retorne erro tentando executar o argumento passado no 'run'
%:
	@: