
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

build-macos: ## Builds the executable for MacOS
	$(MAKE) build GOOS=darwin OUT_DIR=$(OUT_DIR)/macos

build-windows: ## Builds the executable for Windows
	$(MAKE) build GOOS=windows OUT_DIR=$(OUT_DIR)/windows APP_NAME=$(APP_NAME).exe

build-linux: ## Builds the executable for Linux
	$(MAKE) build GOOS=linux OUT_DIR=$(OUT_DIR)/linux

package:	isclean update test build-all ## build a distribution package from a committed state
	mkdir -p dist/package/flightvisualizer
	mv dist/{macos,windows,linux} dist/package/flightvisualizer
	cp .env.local.template dist/package/flightvisualizer
	cp .env.local.template dist/package/flightvisualizer/.env.local
	cp -pR artifacts dist/package/flightvisualizer
	cp docs/package-readme.md dist/package/flightvisualizer/README.md
	(cd dist/package; zip -r ../flightvisualizer.zip .)

build-all:	build-macos build-windows build-linux  ## Builds the executable for all supported architectures

build: ## Build the executable using the currently downloaded dependencies and architecture
	mkdir -p $(OUT_DIR)
	$(GO) build -o $(OUT_DIR)/$(APP_NAME) .

run: ## Run the executable passing it the supplied arguments
	$(GO) run ./... $(ARGS)

test: ## Run the tests
	$(GO) test ./...

clean: ## Remove the build and dependency artifact(s); backup user artifact(s) before removing them
	mkdir -p artifacts.bak
	cp artifacts/* artifacts.bak
	git clean -f -x artifacts
	rm -rf vendor $(OUT_DIR) $(TMP_DIR)

isclean: ## this target fails if there are uncommitted files
	git diff --exit-code
	git diff --exit-code --staged
