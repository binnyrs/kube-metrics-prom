.PHONY: dep
dep:
	@dep ensure -v

.PHONY: build
build: IMG ?= quay.io/bmcstdio/kube-metrics-prom
build: TAG ?= $(VERSION)
build: dep
	docker build -t $(IMG):$(TAG) .

.PHONY: run
run: KUBECONFIG ?= $(HOME)/.kube/config
run:
	@go run cmd/main.go --kubeconfig=$(KUBECONFIG)
