include Makefile.vars

.PHONY: clean
clean:
	@echo "$(ROOT_REPO)@$(BUILD) clean"
	-rm -Rf $(TMP_DIR)

$(TMP_DIR):
	mkdir -p $(TMP_DIR)

.PHONY: clearcache
clearcache:
	@echo "$(ROOT_REPO)@$(BUILD) clearcache"
	-$(foreach path,$(MODULE_LIST),ls -ld $(BASE_DIR)/$(path)/vendor;)

.PHONY: install.gvm
install.gvm:
	@echo "$(ROOT_REPO)@$(BUILD) install.gvm"
	which gvm || \
		curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer | bash

.PHONY: setup.gvm
setup.gvm:
	@echo "$(ROOT_REPO)@$(BUILD) setup.gvm"
	gvm install go$(GOVERSION) -B
	@echo -e 'Please run:\n `gvm use go$(GOVERSION) --default`'

.PHONY: setup
setup: install.deps
	@echo "$(ROOT_REPO)@$(BUILD) setup"

.PHONY: setup.debug
setup.debug: install.debugdeps
	@echo "$(ROOT_REPO)@$(BUILD) setup.debug"

.PHONY: install.deps
install.deps: $(TMP_DIR)
	@echo "$(ROOT_REPO)@$(BUILD) deps"
	which gotestsum || (\
		cd $(TMP_DIR) && \
		curl -O -L https://github.com/gotestyourself/gotestsum/releases/download/v0.3.2/gotestsum_0.3.2_linux_amd64.tar.gz && \
		tar xf gotestsum_0.3.2_linux_amd64.tar.gz && \
		mv -f gotestsum /usr/local/bin \
	)
	gotestsum --help > /dev/null 2>&1
	which codecov || (\
		cd $(TMP_DIR) && \
		curl -L -o codecov https://codecov.io/bash && \
		chmod a+x codecov && \
		mv -f codecov /usr/local/bin \
	)
	codecov -h > /dev/null 2>&1


.PHONY: install.debugdeps
install.debugdeps:
	@echo "$(ROOT_REPO)@$(BUILD) install.debugdeps"
	which dlv || \
		go get -u github.com/derekparker/delve/cmd/dlv
	dlv version > /dev/null 2>&1

.PHONY: docker
docker.build:
	@echo "$(ROOT_REPO)@$(BUILD) docker"
	docker build --build-arg UID=$(shell id -u) --build-arg GID=$(shell id -g) \
		         -t $(DOCKER_NAME) -t $(DOCKER_NAME):$(VERSION) -f $(DOCKER_FILE) .

.PHONY: docker.bash
docker.bash:
	@echo "$(ROOT_REPO)@$(BUILD) docker.bash"
	docker run --rm --name $(NAME)-bash --entrypoint bash -it -u $(shell id -u):$(shell id -g) \
			   -v $(BASE_DIR):/go/src/$(ROOT_REPO) $(DOCKER_NAME)

docker.%:
	@echo "$(ROOT_REPO)@$(BUILD) docker.$*"
	docker run --rm --name $(NAME)-run -u $(shell id -u):$(shell id -g) \
    		    -v $(BASE_DIR):/go/src/$(ROOT_REPO) $(DOCKER_NAME) $*

include .make/*[!\.vars].makefile
