
OUT_DIR=dist
TMP_DIR=tmp
APP_NAME=fviz
GO=go

help: ## Print this help
	@fgrep -h "##" $(MAKEFILE_LIST) | grep -v fgrep | sed 's/\(.*\):.*## \(.*\)/\1 - \2/' | sort

update: ## Update & organize the module dependency tree
	$(GO) mod tidy
	$(GO) mod download
	$(GO) mod vendor

build-darwin: ## Builds the executable for MacOS
	$(MAKE) build GOOS=darwin APP_NAME=$(APP_NAME)_macos

build-windows: ## Builds the executable for Windows
	$(MAKE) build GOOS=windows APP_NAME=$(APP_NAME)_windows

build-linux: ## Builds the executable for Linux
	$(MAKE) build GOOS=linux APP_NAME=$(APP_NAME)_linux

build-all:	build-darwin build-windows build-linux  ## Builds the executable for all supported architectures

build: ## Build the executable using the currently downloaded dependencies and architecture
	$(GO) build -o $(OUT_DIR)/$(APP_NAME) .

run: ## Run the executable passing it the supplied arguments
	$(GO) run ./... $(ARGS)

test: ## Run the tests
	$(GO) test ./...

clean: ## Remove the built and dependency artifact(s)
	rm -rf vendor $(OUT_DIR) $(TMP_DIR)
