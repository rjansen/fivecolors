PROJECT_ID                            ?= project-id
SERVICE_ACCOUNT                       ?= service-account
SERVICE_ACCOUNT_DISPLAYNAME           ?= "Service Account"
SERVICE_ACCOUNT_KEY                   ?= $(HOME)/.hellper/$(SERVICE_ACCOUNT).json
SERVICE_ACCOUNT_EMAIL                 := $(SERVICE_ACCOUNT)@$(PROJECT_ID).iam.gserviceaccount.com
export GOOGLE_APPLICATION_CREDENTIALS ?= $(SERVICE_ACCOUNT_KEY)