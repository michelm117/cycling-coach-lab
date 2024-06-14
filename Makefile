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
	@npx tailwindcss -o assets/public/styles.css --minify

generate-mocks:
	@echo "Generating mocks..."
	@./scripts/generate-mocks.sh

## run: run local project
run: generate
	@echo "Running project..."
	@go run cmd/main.go


## reset database migrations
.PHONY: reset-db
reset-db:
	@echo "Resetting database..."
	@migrate -path ./migrations/ -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" down
	@migrate -path ./migrations/ -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up

## air: run project with air
.PHONY: air
air:
	@docker compose --env-file=.env -f deployments/docker-compose.dev.yml up -d
	@air

.Phony: templ
templ:
	@templ generate -watch -proxy=http://localhost:8080

tailwind:
	@npx tailwindcss -c tailwind.config.js -o assets/public/styles.css --watch

start:
	@echo "Run the following commands in separate terminals to start the project and hot reload the frontend"
	@echo "make air"
	@echo "make tailwind"
	@echo "make templ"

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

.PHONY: test-coverage
test-coverage: generate
	@echo "Running tests with coverage..."
	@go test -coverprofile ./tmp/cover.out ./...
	@go tool cover -html=./tmp/cover.out

.PHONY: update-snapshots
update-snapshots:
	UPDATE_SNAPS=true go test ./...

.PHONY: cypress
cypress:
	@echo "Running cypress tests..."
	@CYPRESS_APP_VERSION=1.2.0 npx cypress run

.PHONY: cypress-open
cypress-open:
	@echo "Running cypress tests..."
	@CYPRESS_APP_VERSION=1.2.0 npx cypress open

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

