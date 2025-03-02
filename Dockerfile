# Build
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
EXPOSE 9111
ARG APP_VERSION=dev
RUN go build -ldflags="-X 'content-alchemist/config.APP_VERSION=${APP_VERSION}'" -o content-alchemist ./cmd/server/main.go

# Runtime
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/content-alchemist .
COPY .env /app/.env
CMD ["./content-alchemist"]