.PHONY: dev test cover tidy fmt

.DEFAULT: help
help:
	@echo "make dev"
	@echo "	setup development environment"
	@echo "make run"
	@echo "	run app"
	@echo "make test"
	@echo "	run go test"
	@echo "make cover"
	@echo "	run go test with -cover"
	@echo "make tidy"
	@echo "	run go mod tidy"
	@echo "make fmt"
	@echo "	run gofmt"
	@echo "make pre-commit"
	@echo "	run pre-commit hooks"

check-pre-commit:
 ifeq (, $(shell which pre-commit))
 $(error "No pre-commit in $(PATH), pre-commit (https://pre-commit.com) is required")
 endif

dev: check-pre-commit
	pre-commit install

run:
	go build . && ./air

test:
	go test -v ./...

cover:
	go test -cover -v ./...

tidy:
	go mod tidy -compat=1.17

fmt:
	gofmt -s -w .

pre-commit: check-pre-commit
	pre-commit
