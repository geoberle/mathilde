.PHONY: lint lint-fix format format-check serve ci

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

ci: lint format-check
