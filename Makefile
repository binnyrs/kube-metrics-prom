VERSION ?= v0.2.0-dev

.PHONY: mod
mod:
	@go mod download

.PHONY: build
build: IMG ?= quay.io/bmcstdio/kube-metrics-prom
build: TAG ?= $(VERSION)
build:
	@docker build -t $(IMG):$(TAG) .

.PHONY: run
run: KUBECONFIG ?= $(HOME)/.kube/config
run: mod
run:
	@go run cmd/main.go --kubeconfig=$(KUBECONFIG)
