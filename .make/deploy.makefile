.PHONY: deploy
ifneq ($(filter $(MODULE_NAME),$(FUNCTION_LIST)),)
deploy: vendor
	echo gcloud beta functions deploy $(FUNCTION_NAME) --source $(MODULE_DIR) --runtime $(FUNCTION_RUNTIME) \
	   	--entry-point $(FUNCTION_ENTRYPOINT) --timeout $(FUNCTION_TIMEOUT) --memory $(FUNCTION_MEMORY) \
		--trigger-http --service-account $(SERVICE_ACCOUNT_EMAIL) --set-env-vars $(FUNCTION_ENVIRONMET)
else
deploy:
	$(error "err_nondeployable_module: module_name=$(MODULE_NAME) avaiable_list=[$(FUNCTION_LIST)]")
endif
