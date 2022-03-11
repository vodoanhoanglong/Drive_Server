.EXPORT_ALL_VARIABLES:

REGISTRY ?= gcr.io/gco-jbee-dev
PROJECT ?= jbee
VERSION ?= $(shell date +"%Y%m%d")
TAG ?= $(shell ./scripts/get-version.sh)
GIT_COMMIT ?= $(shell git rev-parse --short HEAD)
ENV_FILE ?= .env
args=$(filter-out $@,$(MAKECMDGOALS))

dev:
	docker-compose -f docker-compose.yaml -f docker-compose.dev.yaml up -d ${SERVICE}

dev-build:
	docker-compose -f docker-compose.yaml -f docker-compose.dev.yaml up -d --build ${SERVICE}

staging:
	docker-compose -f docker-compose.yaml -f docker-compose.staging.yaml up -d --build ${SERVICE}

clean: 
	docker-compose down --remove-orphans -v

dev-clean: clean dev-build

restart:
	docker-compose restart ${SERVICE}

logs:
	docker-compose logs -f ${SERVICE}

down:
	docker-compose -f docker-compose.yaml -f docker-compose.dev.yaml down ${SERVICE}
	
migrate: 
	./scripts/migrate.sh

bootstrap:
	./scripts/bootstrap/bootstrap.sh

# controller
.PHONY: build-controller
## build-controller: build the controller service
build-controller:
	docker build -t $(REGISTRY)/$(PROJECT)-controller:$(VERSION) services/controller

.PHONY: push-controller
## push-controller: push the controller service to registry
push-controller:
	docker push $(REGISTRY)/$(PROJECT)-controller:$(VERSION)

.PHONY: controller
## controller: build and push the controller service to registry
controller: build-controller push-controller

# auth
.PHONY: build-auth
## build-auth: build the auth service
build-auth:
	docker build -t $(REGISTRY)/$(PROJECT)-auth:$(VERSION) \
		--build-arg TAG=$(TAG) --build-arg GIT_COMMIT=$(GIT_COMMIT) \
		-f services/auth/Dockerfile .

.PHONY: push-auth
## push-auth: push the auth service to registry
push-auth:
	docker push $(REGISTRY)/$(PROJECT)-auth:$(VERSION)

.PHONY: auth
## auth: build and push the auth service to registry
auth: build-auth push-auth

# controller-api
.PHONY: build-controller-api
## build-controller-api: build the controller-api service
build-controller-api:
	docker build -t $(REGISTRY)/$(PROJECT)-api:$(VERSION) \
		--build-arg TAG=$(TAG) --build-arg GIT_COMMIT=$(GIT_COMMIT) \
		-f services/$(PROJECT)-api/Dockerfile .

.PHONY: push-controller-api
## push-controller-api: push the controller-api service to registry
push-controller-api:
	docker push $(REGISTRY)/$(PROJECT)-controller-api:$(VERSION)

.PHONY: controller-api
## controller-api: build and push the controller-api service to registry
controller-api: build-controller-api push-controller-api

.PHONY: run-services
## run-services: start docker services required for integration tests
run-services:
	./scripts/dev.sh dc up -d controller

.PHONY: test-integration
## test-integration: run tests that require docker-compose
test-integration: run-services
	# -p=1 to avoid running tests for different packages in parallel
	go test -v -p=1 -tags=integration ./...

.PHONY: go-unit
## go-unit: run go unit tests for all services
go-unit:
	go test ./...

.PHONY: go-fmt
## go-fmt: check formatting of go code
go-fmt:
	./scripts/check_gofmt.sh

GO_MODULES := $(shell find . -not \( -name vendor -prune \) -name go.mod | xargs -n 1 dirname)

.PHONY: go-vet
## go-vet: check go code with go vet
go-vet:
	go vet ./...

.PHONY: go-lint
## go-lint: lint go code
go-lint: go-fmt go-vet

.PHONY: help
## help: prints help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
