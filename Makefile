# ==================================================
# 🐻 gotterns — Go Library Makefile
# ==================================================
# Commands:
#   make dev-install   → Setup dev environment
#   make lint          → Run linters
#   make fmt           → Auto-format code
#   make tidy          → Clean and verify dependencies
#   make check         → Run lint + tests
#   make test          → Run unit tests
#   make test-cov      → Run tests with coverage
#   make clean         → Clean temporary files
# ==================================================

GO            := go
GOCMD         := @$(GO)
GOTEST        := @$(GOCMD) test
TOOL          := @$(GOCMD) tool
GOLANGCI_LINT := @$(TOOL) -modfile=golangci-lint.mod golangci-lint
PKG_DIRS      := ./...

.DEFAULT_GOAL := help

# ==================================================
# 🧰 Setup
# ==================================================
.PHONY: dev-install
dev-install: ## Install development dependencies
	@echo "🔧 Setting up development environment..."
	sudo apt update -qq
	sudo apt install -y pre-commit
	@$(GOCMD) mod download
	pre-commit install --install-hooks
	@echo "✅ Development environment ready!"

# ==================================================
# 🧹 Lint / Format
# ==================================================
.PHONY: lint
lint: ## Run static analysis using golangci-lint
	@echo "🔍 Running linters..."
	@$(GOLANGCI_LINT) run $(PKG_DIRS)
	@echo "✅ Lint check completed!"

.PHONY: fmt
fmt: ## Auto-format code using golangci-lint and go fmt
	@echo "🧹 Formatting code..."
	@$(GOLANGCI_LINT) run --fix $(PKG_DIRS)
	@$(GOCMD) fmt $(PKG_DIRS)
	@echo "✅ Code formatted successfully!"

# ==================================================
# 🧩 Dependencies
# ==================================================
.PHONY: tidy
tidy: ## Organize and verify Go dependencies
	@echo "🧩 Cleaning and organizing dependencies..."
	@$(GOCMD) mod tidy
	@$(GOCMD) mod verify
	@echo "✅ Dependencies tidy!"

# ==================================================
# 🧪 Tests
# ==================================================
.PHONY: test
test: ## Run all tests
	@echo "🧪 Running tests..."
	@$(GOTEST) -v $(PKG_DIRS)
	@echo "✅ Tests passed!"

.PHONY: test-cov
test-cov: ## Run tests with coverage report
	@echo "📊 Running tests with coverage..."
	@$(GOTEST) -v -coverprofile=coverage.txt $(PKG_DIRS)
	@$(GO) tool cover -func=coverage.txt | grep total
	@echo "✅ Coverage report generated!"

# ==================================================
# 🧽 Clean
# ==================================================
.PHONY: clean
clean: ## Remove temporary files
	@echo "🧽 Cleaning temporary files..."
	@rm -f coverage.txt
	@echo "✅ Clean complete!"

# ==================================================
# ✅ All checks
# ==================================================
.PHONY: check
check: fmt lint test ## Run full validation pipeline (format, lint, test)
	@echo "✅ All checks passed!"

# ==================================================
# 🆘 Help
# ==================================================
.PHONY: help all
help:
	@echo "📘 Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
