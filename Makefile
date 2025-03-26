.PHONY: all dev default install run fmt lint test coverage benchmark build release clean

# Directories and target
BIN=bin/
DIST=dist/
SRC=$(shell find . -name "*.go")
TARGET=$(BIN)/go-httpserver

# Dependency checks
ifeq (, $(shell which golangci-lint))
	$(warning "could not find golangci-lint in $(PATH), \
	run: curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh")
endif

ifeq (, $(shell which goreleaser))
	$(warning "could not find goreleaser in $(PATH), \
	run: go install github.com/goreleaser/goreleaser/v2@latest")
endif

# Default target
default: dev

# Dev target
dev: install fmt lint coverage benchmark

# Meta target
all: dev build

# Commands
install:
	$(info 📥 DOWNLOADING DEPENDENCIES...)
	go get -v ./...

run: build
	$(info ⚙️ RUNNING...)
	@$(TARGET)

fmt:
	$(info ✨ CHECKING CODE FORMATTING...)
	@test -z $(shell gofmt -l $(SRC)) || (gofmt -d $(SRC); exit 1)

lint:
	$(info 🔍 RUNNING LINT TOOLS...)
	golangci-lint run --config .golangci.yaml

test: install
	$(info 🧪 RUNNING TESTS...)
	go test -v ./... -cover

coverage: install
	$(info ✅ TESTING & GENERATING COVERAGE REPORT...)
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html

benchmark: install
	$(info 🚀 RUNNING BENCHMARKS...)
	go test -bench=.


build: install
	$(info 🏗️ BUILDING THE PROJECT...)
	@if [ -e "$(TARGET)" ]; then rm -rf "$(TARGET)"; fi
	@mkdir -p $(BIN)
	@go build -o $(TARGET)

release: fmt lint test benchmark
	$(info 📦 CREATING A NEW RELEASE...)
	goreleaser release

clean:
	$(info 🧹 CLEANING UP...)
	rm -rf $(BIN)
	rm -rf $(DIST)
