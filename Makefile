BUILD_DIR ?= build
BUILD_PACKAGE ?= ./cmd/main.go
REPO_NAME=github.com/bogatyr285/sumcalc-testtask

BINARY_NAME ?= doer-api
VERSION ?= $(shell git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD)
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE ?= $(shell date +%FT%T%z)
LDFLAGS += -s -w -X ${REPO_NAME}/internal/buildinfo.version=${VERSION} -X ${REPO_NAME}/internal/buildinfo.commitHash=${COMMIT_HASH} -X ${REPO_NAME}/internal/buildinfo.buildDate=${BUILD_DATE}

.PHONY: build
build:
	go build ${GOARGS} -tags "${GOTAGS}" -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/${BINARY_NAME} ${BUILD_PACKAGE}

serve:
	go run ./cmd/main.go serve --config config.yaml

test: 
	go test -coverprofile cover.out -race -v ./...

opencoverage:
	go tool cover -html=cover.out
	
jwtkeys:
	openssl ecparam -name prime256v1 -genkey -noout -out private.ec.key
	openssl ec -in private.ec.key -pubout -out public.pem