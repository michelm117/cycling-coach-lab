# FROM golang:1.22 AS build-stage
#
# WORKDIR /app
#
# COPY main /entrypoint
# COPY assets/ /assets
#
# ENV PORT=80
# EXPOSE 80
#
# ENTRYPOINT ["/entrypoint"]

# Build.
FROM golang:1.22 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y nodejs \
    npm
RUN npm install -g tailwindcss
RUN npm install -D daisyui@latest
COPY . /app
RUN templ generate
RUN tailwindcss -o assets/styles.css --minify
RUN CGO_ENABLED=0 GOOS=linux go build -o /entrypoint /app/cmd/main.go

# Deploy.
FROM gcr.io/distroless/static-debian11 AS release-stage
WORKDIR /
COPY --from=build-stage /entrypoint /app/entrypoint
COPY --from=build-stage /app/assets /app/assets
COPY --from=build-stage /app/migrations /app/migrations
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/entrypoint"]

