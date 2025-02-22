# Build
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
EXPOSE 9111
RUN go build -o chappie_server ./cmd/server/main.go

ARG APP_VERSION=dev
RUN go build -ldflags="-X 'chappie_server/config.APP_VERSION=${APP_VERSION}'" -o chappie_server ./cmd/server/main.go

# Runtime
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/chappie_server .
COPY .env /app/.env
CMD ["./chappie_server"]