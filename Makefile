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
	@templ generate
	@go mod tidy
	@npm install


## generate: generate static files
.PHONY: generate
generate:
	@echo "Generating static files..."
	@templ generate 
	@npx tailwindcss -o assets/styles.css --minify


## run: run local project
run: generate
	@echo "Running project..."
	@go run cmd/main.go


## start: build and run project with hot reload
.PHONY: start
start:
	@docker compose --env-file=.env -f deployments/docker-compose.dev.yml up -d
	@air

tailwind:
	@npx tailwindcss -c tailwind.config.js -o assets/styles.css --watch

## update: update project dependencies
.PHONY: update
update:
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@npm update


## test: run unit tests
.PHONY: test
test: generate
	@echo "Running tests..."
	@go test -race -cover ./...


## docker-build: build project into a docker container image
.PHONY: docker-build
docker-build:
	@echo "Building docker image..."
	docker build --no-cache . -t unrealwombat/cycling-coach-lab:latest


## docker-run: run project in a container
.PHONY: docker-run
docker-run:
	@echo "Running docker container..."
	@docker compose -f deployments/docker-compose.prod.yml down --remove-orphans
	@docker compose --env-file=.env -f deployments/docker-compose.prod.yml up

