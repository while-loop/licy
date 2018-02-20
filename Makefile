TARGET := licy
VERSION := 0.0.1
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
SHELL := /usr/bin/env bash

.DEFAULT_GOAL := build

all: update-licenses generate install

clean:
	@echo "[clean] cleaning licy and removing binary files"
	@rm -f ${TARGET}
	@rm -rf bin/ _licenses/

update-licenses:
	@echo "[update-licenses] pulling latest licenses from choosealicense.com"
	@go run scripts/update_licenses.go

generate: update-licenses
	@echo "[generate] generate license files to go files"
	@go generate ./...

build:
	@echo "[build] building go binary"
	@go build \
		-ldflags "-s -w \
		-X main.Version=${VERSION} \
		-X main.Commit=${COMMIT} \
		-X main.Name=${TARGET} \
		-X main.BuildDate=${BUILD_DATE}" \
		-o ${TARGET}

install: build
	@echo "[install] installing go binary to GOPATH"
	@mv ${TARGET} ${GOPATH}/bin/${TARGET}

release:
	@for os in windows darwin linux ; do \
		for arch in amd64 386 ; do \
			GOOS=$$os GOARCH=$$arch go build \
							-ldflags "-s -w \
							-X main.Version=${VERSION} \
							-X main.Commit=${COMMIT} \
							-X main.BuildDate=${BUILD_DATE}" \
							-o bin/${TARGET}-$$os-$$arch; \
			file bin/${TARGET}-$$os-$$arch; \
			if [ "$$os" == "windows" ]; then \
				mv bin/${TARGET}-$$os-$$arch bin/${TARGET}-$$os-$$arch.exe; \
			fi \
		done \
	done
	@echo
	@echo "[release] generating sha256sums"
	@find bin -type f -print0 | xargs -0 sha256sum > bin/sha256sums.txt
	@echo "[release] uploading binaries to GitHub. Tag: v${VERSION}"
	@go run scripts/upload_binaries.go bin ${VERSION}
