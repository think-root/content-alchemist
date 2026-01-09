# Build
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
EXPOSE 9111
ARG APP_VERSION=dev
RUN apk add --no-cache build-base
RUN go mod tidy && \
    CGO_ENABLED=1 go build -ldflags="-X 'content-alchemist/config.APP_VERSION=${APP_VERSION}'" -o content-alchemist ./cmd/server/main.go

# Runtime
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/content-alchemist .
COPY .env /app/.env
CMD ["./content-alchemist"]
