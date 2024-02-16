# Build.
FROM golang:1.22 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -o main /app/cmd/main.go

# Deploy.
FROM gcr.io/distroless/static-debian11 AS release-stage
WORKDIR /
COPY --from=build-stage /app/main /entrypoint
COPY --from=build-stage /app/assets /assets
ENV PORT=80
EXPOSE 80
USER nonroot:nonroot
ENTRYPOINT ["/entrypoint"]
