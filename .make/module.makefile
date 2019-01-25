.PHONY: ci
ci: vet coverage.text bench

.PHONY: build
build: $(TMP_DIR)
	@echo "$(REPO)@$(BUILD) build"
	cd $(MODULE_DIR) && go build -ldflags '-X main.version=$(BUILD)' -o $(MODULE_BIN) $(REPO)

.PHONY: run
run: build
	@echo "$(REPO)@$(BUILD) run"
	$(MODULE_BIN)

.PHONY: vendor
vendor:
	@echo "$(REPO)@$(BUILD) vendor"
	cd $(MODULE_DIR) && go mod verify && go mod vendor

.PHONY: debug
debug:
	@echo "$(REPO)@$(BUILD) debug"
	cd $(MODULE_DIR) dlv debug $(REPO)

.PHONY: debugtest
debugtest:
	@echo "$(REPO)@$(BUILD) debugtest"
	cd $(MODULE_DIR) && dlv test $(DEBUG_PKG) -- -test.run $(TESTS)

.PHONY: vet
vet:
	@echo "$(REPO)@$(BUILD) vet"
	cd $(MODULE_DIR) && go vet $(TEST_PKGS)

.PHONY: test
test:
	@echo "$(REPO)@$(BUILD) test"
	cd $(MODULE_DIR) && gotestsum -f short-verbose -- -v -race -run $(TESTS) $(TEST_PKGS)

.PHONY: itest
itest:
	@echo "$(REPO)@$(BUILD) itest"
	cd $(MODULE_DIR) && gotestsum -f short-verbose -- -tags=integration -v -race -run $(TESTS) $(TEST_PKGS)

.PHONY: bench
bench:
	@echo "$(REPO)@$(BUILD) bench"
	cd $(MODULE_DIR) && gotestsum -f short-verbose -- -bench=. -run="^$$" -benchmem $(TEST_PKGS)

.PHONY: coverage
coverage: $(TMP_DIR)
	@echo "$(REPO)@$(BUILD) coverage"
	ls -ld $(TMP_DIR)
	@touch $(COVERAGE_FILE)
	cd $(MODULE_DIR) && gotestsum -f short-verbose -- -tags=integration -v -run $(TESTS) \
			  -covermode=atomic -coverpkg=$(PKGS) -coverprofile=$(COVERAGE_FILE) $(TEST_PKGS)

.PHONY: coverage.text
coverage.text: coverage
	@echo "$(REPO)@$(BUILD) coverage.text"
	cd $(MODULE_DIR) && go tool cover -func=$(COVERAGE_FILE)

.PHONY: coverage.html
coverage.html: coverage
	@echo "$(REPO)@$(BUILD) coverage.html"
	cd $(MODULE_DIR) && go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@open $(COVERAGE_HTML) || google-chrome $(COVERAGE_HTML) || google-chrome-stable $(COVERAGE_HTML)

.PHONY: coverage.push
coverage.push:
	@echo "$(REPO) coverage.push"
	@#download codecov script and push report with oneline cmd
	@#curl -sL https://codecov.io/bash | bash -s - -f $(COVERAGE_FILE)$(if $(CODECOV_TOKEN), -t $(CODECOV_TOKEN),)
	@codecov -f $(COVERAGE_FILE)$(if $(CODECOV_TOKEN), -t $(CODECOV_TOKEN),)
