PLATFORM ?= linux/arm64,linux/amd64
ENV ?= development
RACE ?= 0
GOPATH ?= $(HOME)/go
APP_NAME ?= vokapi
VERSION ?= latest

EXTENSION ?=

ifeq ($(OS),Windows_NT)
	EXTENSION := .exe
endif

.PHONY: run
run:
	@CXX=g++ CC=gcc go run ./main.go

.PHONY: build
build:
ifeq ($(ENV),production)
	@CGO_ENABLED=0 CXX=g++ CC=gcc go build -ldflags="-s -w" -o ./bin/$(APP_NAME)$(EXTENSION) ./main.go
else ifeq ($(ENV),development)
	@CXX=g++ CC=gcc go build -o ./bin/$(APP_NAME)$(EXTENSION) -gcflags "all=-N -l" ./main.go
else
	@echo "Target ${ENV} is not supported"
endif

.PHONY: test
test:
ifeq ($(RACE), 1)
	@CC=gcc CXX=g++ go test ./... -race -covermode=atomic -coverprofile=coverage.txt -timeout 5m
else
	@CC=gcc CXX=g++ go test ./... -covermode=atomic -coverprofile=coverage.txt -timeout 1m
endif

.PHONY: buildx
buildx:
	@docker buildx build --target production --build-arg VERSION="$(VERSION)" --platform "$(PLATFORM)" -t "brossquad/vokapi:$(VERSION)" --file ./docker/Dockerfile .

.PHONY: tidy
tidy:
	@rm -f go.sum
	@go mod tidy

.PHONY: clean
clean:
	@rm -rf ./bin

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: fmt
fmt:
	@gofumpt -l -w .

.PHONY: gosec
gosec:
	@gosec ./...

.PHONY: update
update:
	go get -u

.PHONY: vet
vet:
	@go vet ./...

.PHONY: zip
zip:
	@zip -r vokapi.zip bin/ systemd.service
