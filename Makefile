SHELL = /bin/bash -e -o pipefail
VERSION=0.1.0
GIT_COMMIT=$(git rev-parse --short HEAD)
LDFLAGS=-ldflags "-s -w -X main.Version=${VERSION} -X main.GitCommit=${GIT_COMMIT}"
BUILD_DIR ?= .build
SOURCE_FILES ?= ./...
TEST_PATTERN ?= .
TEST_OPTIONS ?=
REPORTS_DIR ?= .reports

unit-test:
	@go test $(TEST_OPTIONS) -v \
		-short \
		-failfast \
		-race \
		-run $(TEST_PATTERN)

test:
	@[ -d ${REPORTS_DIR} ] || mkdir -p ${REPORTS_DIR}
	@go test $(TEST_OPTIONS) \
		-failfast \
		-race \
		-coverpkg=./... \
		-covermode=atomic \
		-coverprofile=$(REPORTS_DIR)/coverage.txt $(SOURCE_FILES) \
		-run $(TEST_PATTERN) \
		-bench=$(TEST_PATTERN) \
		-timeout=2m

lint: ## Run lint
	golangci-lint run --timeout 5m

bench:
	@[ -d ${REPORTS_DIR} ] || mkdir -p ${REPORTS_DIR}
	@rm -rf ${REPORTS_DIR}/*
	@go test -run BenchmarkValidatorRequired \
		-bench ValidatorRequired \
		-benchtime 5s \
		-benchmem \
		-memprofile "$(REPORTS_DIR)/benchmem.out"
	@mv "validator.go.test" "${REPORTS_DIR}/bench.test"

bench-prof:
	@go tool pprof -alloc_space "${REPORTS_DIR}/bench.test" "$(REPORTS_DIR)/benchmem.out"

bench-prof-http:
	@go tool pprof -http localhost:6061 -alloc_space "${REPORTS_DIR}/bench.test" "$(REPORTS_DIR)/benchmem.out"
