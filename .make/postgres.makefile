.PHONY: postgres.deps
postgres.deps:
	which migrate || \
		GO111MODULE=off go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate
	migrate -help > /dev/null 2>&1

.PHONY: postgres.run
postgres.run:
	@echo "$(REPO) postgres.run"
	docker run --rm -d --name postgres-run --net=host \
		-v "$(POSTGRES_SCRIPTS_DIR):/docker-entrypoint-initdb.d" postgres:11.1
	@sleep 5 #wait until database is ready

.PHONY: postgres.kill
postgres.kill:
	@echo "$(REPO) postgres.kill"
	docker kill postgres-run

.PHONY: postgres.setup
postgres.setup:
	@echo "$(REPO) postgres.setup"
	sleep 5
	cat $(POSTGRES_SCRIPTS_DIR)/* | \
		psql -h $(POSTGRES_HOST) -U $(POSTGRES_USER) -d $(POSTGRES_DATABASE) -1 -f -

.PHONY: postgres.psql
postgres.psql:
	@echo "$(REPO) postgres.psql"
	docker run --rm -it --name psql-run --net=host postgres:11.1 \
		psql -h $(POSTGRES_HOST) -U $(POSTGRES_USER) -d $(POSTGRES_DATABASE)

.PHONY: postgres.migrations.run
postgres.migrations.run: postgres.run postgres.migrations.up

postgres.migrations.%:
	migrate -source file://$(POSTGRES_MIGRATIONS_DIR) -database $(SQL_DSN) $*
