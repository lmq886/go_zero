GO_PATH := $(shell go env GOPATH)
GO_BIN := $(GO_PATH)/bin
GOCTL := $(GO_BIN)/goctl

.PHONY: all
all: build

.PHONY: build
build:
	go build -o bin/admin-api cmd/admin-api/main.go

.PHONY: run
run:
	go run cmd/admin-api/main.go -f etc/admin-api.yaml

.PHONY: init-db
init-db:
	psql -U postgres -d admin -f scripts/init.sql

.PHONY: gen-api
gen-api:
	@if [ ! -f $(GOCTL) ]; then \
		echo "Installing goctl..."; \
		go install github.com/zeromicro/go-zero/tools/goctl@latest; \
	fi
	$(GOCTL) api go -api api/admin.api -dir api

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -rf bin/
	rm -rf uploads/

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...
