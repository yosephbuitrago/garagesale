SHELL=/bin/bash

# ==============================================================================
# Testing running system

# // To generate a private/public key PEM file.
# openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# openssl rsa -pubout -in private.pem -out public.pem
# ./sales-admin genkey

# curl --user "admin@example.com:gophers" http://localhost:3000/v1/users/token/5cf37266-3473-4006-984f-9325122678b7
# export TOKEN="COPY TOKEN STRING FROM LAST CALL"
# curl -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/users

# hey -m GET -c 100 -n 10000 -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/users
# zipkin: http://localhost:9411
# expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,mem:memstats.Alloc"

# ==============================================================================
all: sales-api

run:
	go run app/sales-api/main.go

runa:
	go run app/admin/admin.go

tidy:
	go mod tidy
	go mod vendor

test:
	go test -v ./...
	# statickcheck

sales-api:
	docker build \
		-f zarf/docker/dockerfile.sales-api \
		-t sales-api-amd64:1.0 \
		--build-arg VCS_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +”%Y-%m-%dT%H:%M:%SZ”` \
		.


# ==============================================================================
# Running from within k8s/dev

kind-up:
	kind create cluster --image kindest/node:v1.19.1 --name yoseph-starter-cluster --config zarf/k8s/dev/kind-config.yaml
	kind load docker-image postgres:13-alpine --name yoseph-starter-cluster

kind-down:
	kind delete cluster --name yoseph-starter-cluster

kind-load:
	kind load docker-image sales-api-amd64:1.0 --name yoseph-starter-cluster
	# kind load docker-image metrics-amd64:1.0 --name yoseph-starter-cluster

kind-services:
	kustomize build zarf/k8s/dev | kubectl apply -f -

kind-sales-api: sales-api
	kind load docker-image sales-api-amd64:1.0 --name yoseph-starter-cluster
	kubectl delete pods -lapp=sales-api

kind-metrics: metrics
	kind load docker-image metrics-amd64:1.0 --name yoseph-starter-cluster
	kubectl delete pods -lapp=sales-api

kind-logs:
	kubectl logs -lapp=sales-api --all-containers=true -f

kind-status:
	kubectl get nodes
	kubectl get pods --watch

kind-status-full:
	kubectl describe pod -lapp=sales-api

kind-shell:
	kubectl exec -it $(shell kubectl get pods | grep sales-api | cut -c1-26) --container app -- /bin/sh

kind-database:
	# ./admin --db-disable-tls=1 migrate
	# ./admin --db-disable-tls=1 seed

kind-delete:
	kustomize build zarf/k8s/dev | kubectl delete -f -
