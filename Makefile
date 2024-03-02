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
	@echo "Initializing project..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/a-h/templ/cmd/templ@latest
	@go mod tidy
	@npm install -D tailwindcss
	@npm install -D daisyui@latest


## generate: generate static files
.PHONY: generate
generate:
	@echo "Generating static files..."
	@tailwindcss -o assets/styles.css --minify
	@templ generate


## run: run local project
run: generate
	@echo "Running project..."
	@go run cmd/main.go


## start: build and run project with hot reload
.PHONY: start
start: generate
	@docker compose --env-file=.env up --build -d --restart=no -f docker-compose.dev.yml
	@air & tailwindcss -o assets/styles.css --minify --watch


## test: run unit tests
.PHONY: test
test: generate
	@echo "Running tests..."
	go test -race -cover ./...


## docker-build: build project into a docker container image
.PHONY: docker-build
docker-build:
	@echo "Building docker image..."
	docker buildx build -t ${name} .
	docker build --no-cache . -t unrealwombat/cycling-coach-lab:latest



## docker-run: run project in a container
.PHONY: docker-run
docker-run:
	@echo "Running docker container..."
	docker compose --env-file=.env -f docker-compose.prod.yml up

