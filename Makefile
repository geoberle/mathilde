.PHONY: lint lint-fix format format-check serve ci build test test-integration test-integration-ci go-vet go-lint

lint:
	npx eslint --no-error-on-unmatched-pattern '**/*.js'

lint-fix:
	npx eslint --fix --no-error-on-unmatched-pattern '**/*.js'

format:
	npx prettier --write .

format-check:
	npx prettier --check .

serve:
	python3 -m http.server 8000

build:
	cd backend && go build ./cmd/mathilde

test:
	cd backend && go test ./...

test-integration:
	cd backend && FIRESTORE_EMULATOR_HOST=localhost:8080 go test ./... -v -count=1

test-integration-ci:
	npx firebase emulators:exec --only firestore --project $(FIRESTORE_PROJECT) \
		'cd backend && FIRESTORE_EMULATOR_HOST=localhost:8080 go test ./... -v -count=1'

FIRESTORE_PROJECT ?= mathilde-61d77

go-vet:
	cd backend && go vet ./...

go-lint: go-vet

ci: lint format-check go-vet test
