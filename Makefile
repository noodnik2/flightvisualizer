
OUT_DIR=dist
TMP_DIR=tmp
APP_NAME=fviz
GO=go

help: ## Print this help
	@fgrep -h "##" $(MAKEFILE_LIST) | grep -v fgrep | sed 's/\(.*\):.*## \(.*\)/\1 - \2/' | sort

update: ## Update & organize the module dependency tree
	$(GO) mod vendor
	$(GO) mod download
	$(GO) mod tidy

build: ## Build the executable using the currently downloaded dependencies
	$(GO) build -o $(OUT_DIR)/$(APP_NAME) .

run: ## Run the executable passing it the supplied arguments
	$(GO) run ./... $(ARGS)

test: ## Run the tests
	$(GO) test ./...

clean: ## Remove the built and dependency artifact(s)
	rm -rf vendor $(OUT_DIR) $(TMP_DIR)
