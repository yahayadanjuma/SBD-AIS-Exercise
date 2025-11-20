
# Stage 1: Build Go 1.23 from source
FROM debian:bookworm-slim AS go-builder
ENV GO_VERSION=1.23.0
RUN apt-get update && apt-get install -y wget git gcc make && \
	wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
	tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
	rm go${GO_VERSION}.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"

# Stage 2: Build your app with Go 1.23
FROM go-builder AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ordersystem main.go

# Stage 3: Minimal runtime image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/ordersystem /app/ordersystem
COPY --from=builder /app/docs /app/docs
EXPOSE 3000
CMD ["/app/ordersystem"]
