.PHONY: setup lint lint-fix format format-check serve ci build test test-integration go-vet go-lint go-lint-fix

GOLANGCI_LINT ?= $(shell go env GOPATH)/bin/golangci-lint

setup:
	@echo "Checking prerequisites..."
	@command -v go >/dev/null 2>&1 || { echo "ERROR: go not found — install from https://go.dev/dl/"; exit 1; }
	@command -v node >/dev/null 2>&1 || { echo "ERROR: node not found — install Node.js 22+"; exit 1; }
	@java -version >/dev/null 2>&1 || { echo "ERROR: java not found — install Java 21+ (needed for Firestore emulator)"; exit 1; }
	@command -v $(GOLANGCI_LINT) >/dev/null 2>&1 || { echo "Installing golangci-lint..."; go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest; }
	@npm ci --loglevel=error --no-fund --no-audit
	cd backend && go mod download
	@echo "Setup complete. Run 'make test-integration' to verify."

lint:
	npx eslint --no-error-on-unmatched-pattern '**/*.js'

lint-fix:
	npx eslint --fix --no-error-on-unmatched-pattern '**/*.js'

format:
	npx prettier --write .
	cd backend && $(GOLANGCI_LINT) fmt ./...

format-check:
	npx prettier --check .
	cd backend && $(GOLANGCI_LINT) fmt --diff ./...

serve:
	python3 -m http.server 8000

build:
	cd backend && go build ./cmd/mathilde

test:
	cd backend && go test ./...

FIRESTORE_PROJECT ?= mathilde-61d77

test-integration: setup
	npx firebase emulators:exec --only firestore --project $(FIRESTORE_PROJECT) \
		'cd backend && FIRESTORE_EMULATOR_HOST=localhost:8080 go test ./... -v -count=1'

go-lint:
	cd backend && $(GOLANGCI_LINT) run -v ./...

go-lint-fix:
	cd backend && $(GOLANGCI_LINT) run -v ./... --fix

go-vet:
	cd backend && go vet ./...

ci: setup lint format-check go-lint test
