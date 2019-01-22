.DEFAULT_GOAL := help

.PHONY: deps
deps: ## Resolve dependencies using vgo.
	GO111MODULE=on go mod download

.PHONY: clean
clean: ## Clean builded binaries.
	rm -rfv build

.PHONY: build
build: ## Build benchmark program and loader program.
	GO111MODULE=on go build -o build/bench cmd/bench/main.go
	GO111MODULE=on go build -o build/loader cmd/loader/main.go

.PHONY: run-bench
run-bench: ## Run benchmark program locally.
	GO111MODULE=on go run cmd/bench/main.go -n 1000

.PHONY: run-loader
run-loader: ## Run loader program locally.
	GO111MODULE=on go run cmd/loader/main.go

.PHONY: test
test: ## Runs unit test.
	GO111MODULE=on go test -cover -race -v ./...

.PHONY: lint
lint: ## Runs static code analysis.
	command -v golint >/dev/null 2>&1 || { go get -u golang.org/x/lint/golint; }
	GO111MODULE=on go list ./... | xargs -L1 golint -set_exit_status

# See "Self-Documented Makefile" article
# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
