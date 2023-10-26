SHELL := /bin/bash

# 固定变量
KIND_CLUSTER := jimmy-cluster
KIND         := kindest/node:v1.27.3
NAMESPACE       := sales-system
APP             := service-pod


run:
	go run main.go

build:
	# 这行，可以更加运行环境的不一样，动态修改程序运行的变量
	go build -ldflags "-X main.build=local"

VERSION := 1.0

all: service

service:
	docker build \
		-f zarf/docker/Dockerfile \
		-t service-arm:1.0 \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

dev-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

dev-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -n=sales-system -o wide --watch

kind-load:
	kind load docker-image service-arm:1.0 --name $(KIND_CLUSTER)

kind-apply:
	cat zarf/k8s/basic/service-pod/basic-service.yaml | kubectl apply -f -

# 同时会清空, namespace
kind-delete:
	cat zarf/k8s/basic/service-pod/basic-service.yaml | kubectl delete -f -


view-images:
	docker exec -it jimmy-cluster-control-plane crictl images


dev-logs:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) --all-containers=true -f --tail=100 --max-log-requests=6 | go run app/tooling/logfmt/main.go -service=$(SERVICE_NAME)

dev-logs-init:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) -f --tail=100 -c init-vault-system
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) -f --tail=100 -c init-vault-loadkeys
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) -f --tail=100 -c init-migrate
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) -f --tail=100 -c init-seed

dev-restart:
	kubectl rollout restart deployment $(APP) --namespace=$(NAMESPACE)

dev-update: all kind-load dev-restart

dev-logs:
	kubectl logs --namespace=$(NAMESPACE) -l app=service --all-containers=true -f --tail=100
