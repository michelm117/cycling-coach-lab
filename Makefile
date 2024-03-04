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
	@cd src; go install github.com/cosmtrek/air@latest
	@cd src; go install github.com/a-h/templ/cmd/templ@latest
	@cd src; go mod tidy
	@npm install -D tailwindcss
	@npm install -D daisyui@latest


## generate: generate static files
.PHONY: generate
generate:
	@echo "Update go mod file"
	@cd src; go mod tidy
	@echo "Generating static files..."
	@npx tailwindcss -o assets/styles.css --minify
	@templ generate -path ./src


## run: run local project
run: generate
	@echo "Running project..."
	@cd src; go run cmd/main.go


## start: build and run project with hot reload
.PHONY: start
start: generate
	@docker compose --env-file=.env -f docker-compose.dev.yml up -d
	@cd src; air & tailwindcss -c ../tailwind.config.js -o ../assets/styles.css --minify --watch


## update: update project dependencies
.PHONY: update
update:
	@echo "Updating dependencies..."
	@cd src; go get -u ./...
	@cd src; go mod tidy
	@npm update


## test: run unit tests
.PHONY: test
test: generate
	@echo "Running tests..."
	@cd src; go test -race -cover ./...


## docker-build: build project into a docker container image
.PHONY: docker-build
docker-build:
	@echo "Building docker image..."
	docker build --no-cache . -t unrealwombat/cycling-coach-lab:latest


## docker-run: run project in a container
.PHONY: docker-run
docker-run:
	@echo "Running docker container..."
	@docker compose -f docker-compose.prod.yml down --remove-orphans
	@docker compose --env-file=.env -f docker-compose.prod.yml up

