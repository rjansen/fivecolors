.PHONY: postgres.setup
postgres.setup:
	which migrate || (\
		cd $(TMP_DIR) && \
		curl -O -L https://github.com/golang-migrate/migrate/releases/download/v4.2.3/migrate.linux-amd64.tar.gz && \
		tar xf migrate.linux-amd64.tar.gz && \
		mv -f migrate.linux-amd64 /usr/local/bin/migrate \
	)
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

.PHONY: postgres.scripts
postgres.scripts:
	@echo "$(REPO) postgres.scripts"
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
