GQLGEN_CONFIG ?= $(BASE_DIR)/gqlgen.yml

.PHONY: gqlgen.init
gqlgen.init:
	go run github.com/99designs/gqlgen -v init

.PHONY: gqlgen.generate
gqlgen.generate:
	go run github.com/99designs/gqlgen -v generate
