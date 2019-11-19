NAMESPACE   ?= fivecolors
MODULE      ?= graphql
APP         ?= asset
VERSION     ?= 0.0.1
BUILD_DIR   ?= $(BASE_DIR)/build
DOCKER_DIR  ?= $(BUILD_DIR)/docker
K8S_DIR     ?= $(BUILD_DIR)/k8s


.PHONY: deploy.default
deploy.default: deploy.build

.PHONY: deploy.build
deploy.build:
	docker build -t e-pedion/fivecolors-$(MODULE):$(VERSION) -t e-pedion/fivecolors-$(MODULE):latest -f $(DOCKER_DIR)/$(MODULE)/Dockerfile .

.PHONY: deploy.install.ingress
deploy.install.ingress:
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/mandatory.yaml
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/cloud-generic.yaml

.PHONY: deploy.namespace
deploy.namespace:
	kubectl create namespace $(NAMESPACE)

.PHONY: deploy.secret.credentials
deploy.secret.credentials:
	kubectl -n $(NAMESPACE) create secret generic $(SECRET_CREDENTIALS) --from-file=key.json=$(SECRET_CREDENTIALS_FILE)

.PHONY: deploy.firestore
deploy.firestore:
	kubectl apply -f $(K8S_DIR)/firestore.yml

.PHONY: deploy.ingress
deploy.ingress:
	kubectl apply -f $(K8S_DIR)/ingress.yml

.PHONY: deploy
deploy: deploy.build
	kubectl apply -f $(K8S_DIR)/$(MODULE).yml 

# bkp
.PHONY: deploy_foward
deploy_foward:
	kubectl port-forward -n $(NAMESPACE) svc/$(MODULE) 8080

.PHONY: deploy_setimg
deploy_setimg:
	# i guess that this command overrides the image value to another
	kubectl create -f $(K8S_DIR)/$(MODULE).yml --set spec.template.spec.containers[0].image=e-pedion/fivecolors-$(MODULE):$(VERSION)

redeploy.%: build
	kubectl set image -n $(NAMESPACE) fivecolors-graphql=e-pedion/fivecolors-graphql:deploy-$*

release.%:
	docker build -t gcr.io/cdp-development/graphql-ftw:deploy-$* .
	docker push gcr.io/cdp-development/graphql-ftw:deploy-$*
