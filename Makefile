TEST_PKGS := 

install:
	go build farm.e-pedion.com/repo/fivecolors

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

run: install
	./fivecolors