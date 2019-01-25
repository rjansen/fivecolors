MODULE_NAME        ?= core
MODULE_LIST        := core server asset
MODULE_PATH        := $(filter $(MODULE_NAME),$(MODULE_LIST))
MODULE             ?= $(if $(MODULE_PATH),$(MODULE_PATH),$(firstword $(MODULE_LIST)))
REPO               := $(ROOT_REPO)/$(MODULE)
MODULE_DIR         := $(BASE_DIR)/$(MODULE)
MODULE_BIN         := $(TMP_DIR)/$(MODULE_NAME)
GOVERSION          := 1.11.4
export GO111MODULE := on
PKGS               := ./...
TEST_PKGS          := $(if $(TEST_PKGS),$(addprefix $(REPO)/,$(TEST_PKGS)),$(PKGS))
TESTS              ?= .
DEBUG_PKG          := $(if $(filter $(TEST_PKGS),$(PKGS)),$(REPO),$(firstword $(TEST_PKGS)))
COVERAGE_FILE      := $(TMP_DIR)/$(MODULE_NAME).coverage
COVERAGE_HTML      := $(TMP_DIR)/$(MODULE_NAME).coverage.html