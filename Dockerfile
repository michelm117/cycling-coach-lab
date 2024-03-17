# Build.
FROM golang:1.22 AS build-stage
COPY . /app
## Compile the Go code.
WORKDIR /app/src
RUN go mod download
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate 
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/entrypoint /app/src/cmd/main.go
## Generate the CSS.
WORKDIR /app
RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y nodejs \
    npm
RUN npm install
RUN npx tailwindcss -o /app/src/assets/styles.css --minify


# Deploy.
FROM gcr.io/distroless/static-debian11 AS release-stage
WORKDIR /app
COPY --from=build-stage /app/entrypoint /app/entrypoint
COPY --from=build-stage /app/src/assets /app/src/assets
COPY --from=build-stage /app/src/migrations /app/src/migrations
ENV ENV=production
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/entrypoint"]

