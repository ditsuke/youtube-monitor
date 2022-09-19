SHELL := /bin/bash

gen:
	if [ -d "model" ]; then rm -rf gen; fi
	go run ./cmd/generate/main.go

.PHONY: test
