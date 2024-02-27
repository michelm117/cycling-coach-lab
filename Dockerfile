# FROM golang:1.22 AS build-stages
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
RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y nodejs \
    npm
RUN npm install -g tailwindcss
RUN npm install -D daisyui@latest
RUN tailwindcss -o assets/styles.css --minify
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate
COPY . /app
COPY go.mod go.sum ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /entrypoint /app/cmd/main.go

# Deploy.
FROM gcr.io/distroless/static-debian11 AS release-stage
WORKDIR /
COPY --from=build-stage /entrypoint /entrypoint
COPY --from=build-stage /app/assets /assets
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/entrypoint"]
