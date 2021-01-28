.PHONY: build
build:
	go run -v ./cmd/httpserver/main.go

.DEFAULT_GOAL := build