PACKAGES := $(shell go list ./...)
name := $(shell basename ${PWD})

all: help

.PHONY: help
help: Makefile
	@echo
	@echo " Choose a make command to run"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## init: initialize project (make init module=github.com/user/project)
.PHONY: init
init:
	@go install github.com/cosmtrek/air@latest
	@npm install -g tailwindcss
	@go mod download

## generate: generate static files
.PHONY: generate
generate:
	@tailwindcss -o assets/styles.css --minify
	@templ generate


## run: run local project
run: generate
	@go run cmd/main.go

## start: build and run project with hot reload
.PHONY: start
start: generate
	@air & tailwindcss -o assets/styles.css --minify --watch

## test: run unit tests
.PHONY: test
test:
	go test -race -cover $(PACKAGES)

## docker-build: build project into a docker container image
.PHONY: docker-build
docker-build: test
	GOPROXY=direct docker buildx build -t ${name} .

## docker-run: run project in a container
.PHONY: docker-run
docker-run:
	docker run -it --rm -p 80:80 ${name}

