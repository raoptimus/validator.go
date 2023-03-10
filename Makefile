SHELL = /bin/bash -e -o pipefail
VERSION=0.1.0
GIT_COMMIT=$(git rev-parse --short HEAD)
LDFLAGS=-ldflags "-s -w -X main.Version=${VERSION} -X main.GitCommit=${GIT_COMMIT}"
BUILD_DIR ?= .build
SOURCE_FILES ?= ./...
TEST_PATTERN ?= .
TEST_OPTIONS ?=
COVERAGE_DIR ?= .reports

unit-test:
	@go test $(TEST_OPTIONS) -v \
		-short \
		-failfast \
		-race \
		-run $(TEST_PATTERN)

test:
	@[ -d ${COVERAGE_DIR} ] || mkdir -p ${COVERAGE_DIR}
	@go test $(TEST_OPTIONS) \
		-failfast \
		-race \
		-coverpkg=./... \
		-covermode=atomic \
		-coverprofile=$(COVERAGE_DIR)/coverage.txt $(SOURCE_FILES) \
		-run $(TEST_PATTERN) \
		-bench=$(TEST_PATTERN) \
		-timeout=2m

bench:
	@cd rule && go test -run $(TEST_PATTERN) \
		-bench ValidatorUrl \
		-benchtime 3s \
		-benchmem -memprofile "../$(COVERAGE_DIR)/mem.out"
	@go tool pprof -alloc_space rule/rule.test "$(COVERAGE_DIR)/mem.out"
