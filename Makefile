## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
# =========================================================================
# Dev Tools (Cool but not all necessary tools)
# Database Access
# dblab --host 0.0.0.0 --user postgres --db postgres --pass postgres --ssl disable --port 5432 --driver postgres
# https://github.com/danvergara/dblab

# Load Testing:
# hey -m GET -c 100 -n 10000 http://localhost:3000/app/list


## dev.setup.mac: Install commonly used tools for the development process.
dev.setup.mac:
	brew update
	brew list go || brew install go
	brew list kind || brew install kind
	brew list kubectl || brew install kubectl
	brew list kustomize || brew install kustomize
	brew list jq || brew install jq
# brew list hey || brew install hey
# =========================================================================
# Building Containers

## VERSION: Version of application to set.
VERSION := 1.0.0

## image: run appstore-api to build docker image.
image: appstore-api

## appstore-api: Build docker image for api from zarf/docker/dockerfile.appstore-api
appstore-api:
	docker build \
		-f zarf/docker/dockerfile.appstore-api \
		-t appstore-api-arm64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u + "%Y-%m-%dT%H:%M:%SZ"` \
		.

#=====================================================================
# KIND image release info at project: github.com/kubernetes-sigs/kind/releases/tag/[your version of kind]
# Running from within k8s/kind

## KIND_CLUSTER: Kind cluster name default 'gopherhelx-cluster'
KIND_CLUSTER := gopherhelx-cluster 

## kind-up: Start new cluster, set namespace 'appstore-system', patching changes.
kind-up:
	kind create cluster \
		--image kindest/node:v1.25.3@sha256:f52781bc0d7a19fb6c405c2af83abfeb311f130707a0e219175677e366cc45d1 \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yml
	kubectl config set-context --current --namespace=appstore-system

## kind-default-CRB: Create cluster rolebinding for api to manage deployments/services.
# Only needed to run once when first creating cluster.
# Still honing this so for hacking, it is cluster-admin roles (VERY INSECURE!!!)
kind-default-CRB:
	kubectl create clusterrolebinding serviceaccounts-cluster-admin \
	--clusterrole=cluster-admin \
	--group=system:serviceaccounts

cd-example:
	cd ~ 
	echo "HELLO"

# kubectl create clusterrole \
# appstore-manager-role \
# --verb=get,list,watch,update,delete,create,patch \
# --resource=*.apps,deployment.apps,pods,pods/status,deployments,deployments/status,services,services/status,replicasets
# kubectl create clusterrolebinding \
# appstore-manager-binding \
# --namespace=appstore-system \
# --clusterrole=appstore-manager-role \
# --serviceaccount=default:default

# To verify serviceaccount and clusterrolebinding
# kubectl get rolebindings,clusterrolebindings \
  --all-namespaces  \
  -o custom-columns='KIND:kind,NAMESPACE:metadata.namespace,NAME:metadata.name,SERVICE_ACCOUNTS:subjects[?(@.kind=="ServiceAccount")].name' | grep "default"


# NOTE: This role and role binding are used by rest.InClusterConfig() in business/k8s package
# in order to authenticate the appstore pod to the k8s api server.

## kind-down: Delete the kind cluster
kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

## kind-load: Load image into cluster.
kind-load:
	cd zarf/k8s/kind/appstore; kustomize edit set image appstore-api-image=appstore-api-arm64:$(VERSION)
	kind load docker-image appstore-api-arm64:$(VERSION) --name $(KIND_CLUSTER)

# Load db first, wait 120, then apply appstore api
## kind-apply: Apply manifest and kustomize patches into kubernetes.
kind-apply:
	kustomize build zarf/k8s/kind/database | kubectl apply -f -
	kubectl wait --namespace=appstore-system --timeout=120s --for=condition=Available deployment/database-pod
	kustomize build zarf/k8s/kind/appstore | kubectl apply -f -

## kind-status: Get status of nodes, svc, and pods.
kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide 

## kind-logs: Actively watch logging, starting at last 100 logs.
kind-logs:
	kubectl logs -l app=appstore -f --tail=100 

## kind-restart: Rollout and restart new deployment of appstore-api
kind-restart:
	kubectl rollout restart deployment appstore-api


## kind-restart-all: Rollout and restart new deployments api and db.
kind-restart-all:
	kubectl rollout restart deployment appstore-api
	kubectl rollout restart deployment database-pod

## kind-status-appstore: Get status of just api pod
kind-status-appstore:
	kubectl get pods -o wide -w -l app=appstore

## kind-status-db: Get status of the database pod
kind-status-db:
	kubectl get pods -o wide -w -l app=database

## kind-update: Update runs docker build and restarts the deployment with new rollout.
# Use this if you edit application code.
kind-update: image kind-load kind-restart

## kind-update-apply: Runs docker build, load and kustomize patching.
# Use this if you edit zarf/ files.
kind-update-apply: image kind-load kind-apply

## kind-describe: Describes the pods with label=appstore
kind-describe:
	kubectl describe pod -l app=appstore

# ==============================================================================
# Administration

## migrate: Run admin data migrations on db.
# migrate:
# 	go run app/tooling/admin/main.go migrate

# ## seed: Will make and seed the db with new information. 
# seed: migrate
# 	go run app/tooling/admin/main.go seed

# delete: 
# 	go run app/tooling/admin/main.go 3

#===================================================================
# Managing go packages by downloading with 'tidy'
# keeping local copy of package dependency with 'vendor'

## tidy: This will use go mod tidy and go mod vendor.
tidy:
	go mod tidy 
	go mod vendor

## audit: Runs formatting and vetting check.
audit:
	@echo 'Formatting code...' go fmt ./...
	@echo 'Vetting code...'
	go vet ./...

