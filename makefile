SHELL := /bin/bash


#====================== Testing System ===========================+
# Access metrics directly (4000) or through the sidecar (3001)
# go install github.com/divan/expvarmon@latest
run-exp:
 	expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"
# expvarmon -ports=":3001" -endpoint="/metrics" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"


# 固定变量
KIND_CLUSTER := jimmy-cluster
KIND         := kindest/node:v1.27.3
NAMESPACE       := sales-system
APP             := sales-pod

# SERVICE_IMAGE   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(VERSION)
SERVICE_IMAGE := salse-api

# RSA Keys
# 	To generate a private/public key PEM file.
# 	$ openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# 	$ openssl rsa -pubout -in private.pem -out public.pem
# 	$ ./sales-admin genkey

admin:
	go run app/tooling/admin/main.go

run:
	go run app/services/sales-api/main.go | go run app/tooling/logfmt/main.go

build:
	# 这行，可以更加运行环境的不一样，动态修改程序运行的变量
	go build -ldflags "-X main.build=local"

VERSION := 1.1

all: sales-api

sales-api:
	docker build \
		-f zarf/docker/sales-api.Dockerfile \
		-t sales-api-image-arm:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

dev-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

# 同时会清空, namespace
dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

dev-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -n=sales-system -o wide --watch

kind-load:
	cd zarf/k8s/kind/sales-pod; kustomize edit set image sales-api=sales-api-amd64:$(SERVICE_IMAGE)

	# FIXME: 这里镜像的名称和镜像的版本都是硬编码，后面有巨大的问题。
	kind load docker-image sales-api-image-arm:1.1 --name $(KIND_CLUSTER)

kind-apply:
	# 以前不用 kustomize 时的命令
	# cat zarf/k8s/basic/services-pod/basic-services.yaml | kubectl apply -f -
	kustomize build zarf/k8s/kind/sales-pod | kubectl apply -f -
	# kubectl rollout status deployment/sales-pod

kind-delete:
	cat zarf/k8s/basic/sales-pod/basic-sales.yaml | kubectl delete -f -

view-images:
	docker exec -it jimmy-cluster-control-plane crictl images

#dev-logs:
#	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) --all-containers=true -f --tail=100 --max-log-requests=6 | go run app/tooling/logfmt/main.go -service=$(SERVICE_NAME)

dev-logs-init:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) -f --tail=100 -c init-vault-system
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) -f --tail=100 -c init-vault-loadkeys
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) -f --tail=100 -c init-migrate
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) -f --tail=100 -c init-seed

dev-restart:
	kubectl rollout restart deployment $(APP) --namespace=$(NAMESPACE)

dev-update: all kind-load dev-restart

dev-update-apply: all kind-load kind-apply

dev-logs:
	kubectl logs --namespace=$(NAMESPACE) -l app=sales --all-containers=true -f --tail=100

kind-up:
	kind create cluster \
		--image kindest/node:v1.25.2@sha256:9be91e9e9cdf116809841fc77ebdb8845443c4c72fe5218f3ae9eb57fdb4bace \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=sales-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-restart:
	kubectl rollout restart deployment sales-pod

kind-update: all kind-load kind-restart


#
# Deploy First
# make admin (取这里的 token 进去）
# For testing a simple query on the system. Don't forget to `make seed` first.
# curl --user "admin@example.com:gophers" http://localhost:3000/v1/users/token
# export TOKEN="COPY TOKEN STRING FROM LAST CALL"
# curl -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/users/1/2
#

# curl -H "" http://localhost:3000/v1/testauth
# curl -il "" http://localhost:3000/v1/testauth
#

test:
	go test -count=1 ./...
	staticcheck -checks=all ./...
	govulncheck ./...

