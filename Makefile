NAME 		:= fivecolors
BIN         := $(NAME)
REPO        := github.com/rjansen/$(NAME)
BUILD       := $(shell git rev-parse --short HEAD)
#VERSION     := $(shell git describe --tags $(shell git rev-list --tags --max-count=1))
MAKEFILE    := $(word $(words $(MAKEFILE_LIST)), $(MAKEFILE_LIST))
BASE_DIR    := $(shell cd $(dir $(MAKEFILE)); pwd)
PKGS        := $(shell go list ./... | grep -v /vendor/)
COVERAGE_FILE   := $(NAME).coverage
COVERAGE_HTML  	:= $(NAME).coverage.html
PKG_COVERAGE   	:= $(NAME).pkg.coverage

ETC_DIR := ./etc
CONF_DIR := $(ETC_DIR)/$(NAME)
CONF_TYPE := yaml
CONF := $(CONF_DIR)/$(NAME).$(CONF_TYPE)

TEST_PKGS := 

ENV ?= local
TEST_PKGS ?= 

.PHONY: default
default: build

.PHONY: local
local: 
	@echo "Set enviroment to local"
	$(eval ENV = "local")

.PHONY: dev
dev: 
	@echo "Set enviroment to dev"
	$(eval ENV = "dev")

.PHONY: prod
prod: 
	@echo "Set enviroment to prod"
	$(eval ENV = "prod")

.PHONY: check_env
check_env:
	@if [ "$(ENV)" == "" ]; then \
	    echo "Env is blank: $(ENV)"; \
	    exit 540; \
	fi

.PHONY: build
build:
	go build $(REPO)

.PHONY: run
run: build
	./$(NAME) -ecf $(CONF)

pkg_data:
	@echo "Add data pkg for tests"
	$(eval TEST_PKGS += "farm.e-pedion.com/repo/fivecolors/data")

pkg_api:
	@echo "Add api pkg for tests"
	$(eval TEST_PKGS += "farm.e-pedion.com/repo/fivecolors/api")

pkg_test: pkg_data pkg_api
	@echo "TEST_PKGS=$(TEST_PKGS)"

test:
	@if [ "$(TEST_PKGS)" == "" ]; then \
	    echo "Build Without TEST_PKGS" ;\
	    go test farm.e-pedion.com/repo/fivecolors/data farm.e-pedion.com/repo/fivecolors/api ;\
	else \
	    echo "Build With TEST_PKGS=$(TEST_PKGS)" ;\
	    go test $(TEST_PKGS) ;\
	fi
