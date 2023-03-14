## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# =========================================================================
# Building Containers

## VERSION: of application to set.
VERSION := 1.0.0

## all: run appstore-api to build docker image.
all: appstore-api

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

## KIND_CLUSTER: Kind cluster name
KIND_CLUSTER := gopherhelx-cluster 

## kind-up: Start new cluster, set context to appstore-system, reference kind for patching.
kind-up:
	kind create cluster \
		--image kindest/node:v1.25.3@sha256:f52781bc0d7a19fb6c405c2af83abfeb311f130707a0e219175677e366cc45d1 \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yml
	kubectl config set-context --current --namespace=appstore-system

## kind-down: Delete the kind cluster
kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

## kind-load: Use kustomize to replace our VERSION from VERSION in makefile. Load image into cluster
kind-load:
	cd zarf/k8s/kind/appstore; kustomize edit set image appstore-api-image=appstore-api-arm64:$(VERSION)
	kind load docker-image appstore-api-arm64:$(VERSION) --name $(KIND_CLUSTER)

# Load db first, wait 120, then apply appstore api
## kind-apply: Apply kustomize build into kubernetes.
kind-apply:
	kustomize build zarf/k8s/kind/database | kubectl apply -f -
	kubectl wait --namespace=database-system --timeout=120s --for=condition=Available deployment/database-pod
	kustomize build zarf/k8s/kind/appstore | kubectl apply -f -

## kind-status: Get status of nodes, svc, and pods in all namespaces.
kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --all-namespaces

## kind-logs: See last 100 logs for appstore
kind-logs:
	kubectl logs -l app=appstore -f --tail=100 

## kind-restart: Rollout and restart new deployment of appstore-api
kind-restart:
	kubectl rollout restart deployment appstore-api

## kind-status-appstore: Get status of just pods
kind-status-appstore:
	kubectl get pods -o wide -w

## kind-status-db: Get status of the database
kind-status-db:
	kubectl get pods -o wide -w --namespace=database-system

## kind-update: Update runs docker build and restarts the deployment with new rollout.
kind-update: all kind-load kind-restart

## kind-update-apply: Runs docker build, load and kustomize patching.
kind-update-apply: all kind-load kind-apply

## kind-describe: Describes the pods with label=appstore
kind-describe:
	kubectl describe pod -l app=appstore

# ==============================================================================
# Administration

migrate:
	go run app/tooling/admin/main.go migrate

seed: migrate
	go run app/tooling/admin/main.go seed

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

