GO_MODULE_NAME     ?= core
GO_MODULE_LIST     := core function cmd/server cmd/asset
GO_MODULE_PATH     := $(or $(filter $(GO_MODULE_NAME),$(GO_MODULE_LIST)),$(firstword $(GO_MODULE_LIST)))
GO_MODULE          := $(REPO)/$(GO_MODULE_PATH)
GO_MODULE_DIR      := $(BASE_DIR)/$(GO_MODULE_PATH)
GO_MODULE_BIN      := $(TMP_DIR)/$(GO_MODULE_PATH)
GO_VERSION         := 1.13.3
export GO111MODULE := on
GO_PKGS            := ./...
GO_TEST_PKGS       := $(if $(GO_TEST_PKGS),$(addprefix $(GO_MODULE)/,$(GO_TEST_PKGS)),$(GO_PKGS))
GO_TESTS           ?= .
GO_TEST_VERBOSITY  ?= $(if $(DEBUG),standard-verbose,short-verbose)
GO_COVERAGE_FILE   := $(TMP_DIR)/$(subst /,_,$(GO_MODULE_PATH)).coverage
GO_COVERAGE_HTML   := $(TMP_DIR)/$(subst /,_,$(GO_MODULE_PATH)).coverage.html

.PHONY: go.clearcache
go.clearcache:
	@echo "$(NAME)@$(BUILD) go.clearcache"
	-rm -f $(GO_COVERAGE_FILE)
	-rm -f $(GO_COVERAGE_HTML)
	-rm -rf vendor

.PHONY: go.clearmodules
go.clearmodules:
	@echo "$(NAME)@$(BUILD) go.clearmodules"
	-rm -Rf $(HOME)/.cache/go-build
	-rm -Rf $(HOME)/gopath/pkg/mod

.PHONY: go.build
go.build: $(TMP_DIR)
	@echo "$(NAME)@$(BUILD) go.build"
	cd $(GO_MODULE_DIR) && go build -ldflags '-X main.version=$(BUILD)' -o $(GO_MODULE_BIN) $(GO_MODULE_REPO)

.PHONY: go.run
go.run: go.build
	@echo "$(NAME)@$(BUILD) go.run"
	$(GO_MODULE_BIN)

.PHONY: go.vendor
go.vendor:
	@echo "$(NAME)@$(BUILD) go.vendor"
	cd $(GO_MODULE_DIR) && go mod verify && go mod vendor

.PHONY: go.vet
go.vet:
	@echo "$(NAME)@$(BUILD) go.vet"
	cd $(GO_MODULE_DIR) && go vet $(GO_TEST_PKGS)

.PHONY: go.lint
go.lint:
	@echo "$(NAME)@$(BUILD) go.lint"
	cd $(GO_MODULE_DIR) && golint $(shell go list $(GO_TEST_PKGS))

.PHONY: go.test
go.test:
	@echo "$(NAME)@$(BUILD) go.test $(GO_TEST_PKGS)"
	cd $(GO_MODULE_DIR) && $(GOTESTSUM) -f $(GO_TEST_VERBOSITY) -- -v -race -count 1 -run $(GO_TESTS) $(GO_TEST_PKGS)

.PHONY: go.itest
go.itest:
	@echo "$(NAME)@$(BUILD) go.itest $(GO_TEST_PKGS)"
	cd $(GO_MODULE_DIR) && $(GOTESTSUM) -f $(GO_TEST_VERBOSITY) -- -tags="integration" -v -race -count 1 -run $(GO_TESTS) $(GO_TEST_PKGS)

.PHONY: go.bench
go.bench:
	@echo "$(NAME)@$(BUILD) go.bench"
	cd $(GO_MODULE_DIR) && $(GOTESTSUM) -f $(TEST_VERBOSITY) -- -bench=. -run="^$$" -benchmem $(GO_TEST_PKGS)

.PHONY: go.coverage
go.coverage: $(TMP_DIR)
	@echo "$(NAME)@$(BUILD) go.coverage $(GO_TEST_PKGS)"
	@touch $(GO_COVERAGE_FILE)
	cd $(GO_MODULE_DIR) && $(GOTESTSUM) -f $(GO_TEST_VERBOSITY) -- -tags="integration" -v -race -count 1 -run $(GO_TESTS) \
		-covermode=atomic -coverpkg=$(GO_PKGS) -coverprofile=$(GO_COVERAGE_FILE) $(GO_TEST_PKGS)

.PHONY: go.coverage.text
go.coverage.text: go.coverage
	@echo "$(NAME)@$(BUILD) go.coverage.text"
	cd $(GO_MODULE_DIR) && go tool cover -func=$(GO_COVERAGE_FILE)

.PHONY: go.coverage.html
go.coverage.html: go.coverage
	@echo "$(NAME)@$(BUILD) go.coverage.html"
	cd $(GO_MODULE_DIR) && go tool cover -html=$(GO_COVERAGE_FILE) -o $(GO_COVERAGE_HTML)
