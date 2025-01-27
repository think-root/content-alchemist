# Build
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o chappie ./cmd/server/main.go

# Runtime
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/chappie .
COPY .env /app/.env
CMD ["./chappie"]