# syntax=docker/dockerfile:1

############################
# 1) Builder stage
############################
FROM golang:1.25 AS builder

WORKDIR /app
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Cache modules first
COPY go.mod go.sum ./
RUN go mod download

# Copy full source and build
COPY . .
RUN echo "=== Building Go application ==="
RUN chmod +x scripts/build-application.sh && scripts/build-application.sh

# Sanity check
RUN test -f /app/ordersystem

############################
# 2) Runtime stage
############################
FROM alpine:3.20 AS run

WORKDIR /app
RUN apk add --no-cache ca-certificates && update-ca-certificates

# Non-root user
RUN adduser -D -H -s /sbin/nologin appuser

# Copy binary from builder
COPY --from=builder /app/ordersystem /app/ordersystem

# Default environment variables (can be overridden via --env-file)
ENV POSTGRES_DB=order \
    POSTGRES_USER=docker \
    POSTGRES_PASSWORD=docker123 \
    POSTGRES_TCP_PORT=5432 \
    DB_HOST=127.0.0.1

# Expose your app’s HTTP port (adjust if different)
EXPOSE 8080

USER appuser
ENTRYPOINT ["/app/ordersystem"]
