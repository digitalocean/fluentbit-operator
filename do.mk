IMAGE_TAG ?= $(shell git rev-parse --short HEAD)

FB_VERSION ?= 2.2.2

ifdef release
	REV = $(shell git rev-list --tags --max-count=1)
	IMAGE_TAG = $(shell git describe --tags $(REV))
endif

$(info using image tag: $(IMAGE_TAG))

.PHONY: image-operator
image-operator:
	docker build -f cmd/manager/Dockerfile -t digitaloceanapps/fluent-bit-operator:$(IMAGE_TAG) .
ifdef latest
	docker tag digitaloceanapps/fluent-bit-operator:$(IMAGE_TAG) digitaloceanapps/fluent-bit-operator:latest
endif

.PHONY: image-push-operator
image-push-operator: image-operator
	docker push digitaloceanapps/fluent-bit-operator:$(IMAGE_TAG)
ifdef latest
	docker push digitaloceanapps/fluent-bit-operator:latest
endif

.PHONY: image-fluentbit
image-fluentbit:
	docker build --platform linux/amd64 --build-arg FLUENTBIT_VERSION=$(FB_VERSION) -f cmd/fluent-bit-watcher/Dockerfile -t digitaloceanapps/fluent-bit:$(FB_VERSION)-$(IMAGE_TAG) .
ifdef latest
	docker tag digitaloceanapps/fluent-bit:$(FB_VERSION)-$(IMAGE_TAG) digitaloceanapps/fluent-bit:latest
endif

.PHONY: image-push-fluentbit
image-push-fluentbit: image-fluentbit
	docker push digitaloceanapps/fluent-bit:$(FB_VERSION)-$(IMAGE_TAG)
ifdef latest
	docker push digitaloceanapps/fluent-bit:latest
endif