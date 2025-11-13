# ---- Build stage: compile the Go orderservice ----
FROM golang:1.25 AS builder

# Work inside /app in the container
WORKDIR /app

# Copy all project files into /app
COPY . .

# Download dependencies and write go.sum
RUN go mod tidy

# Build the orderservice binary as static
RUN go build -a -installsuffix cgo -ldflags="-w -s" -o /app/ordersystem .


# ---- Run stage: small runtime image ----
FROM alpine:3.20

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates libc6-compat

WORKDIR /

# Copy the built binary from the builder stage
COPY --from=builder /app/ordersystem /app/ordersystem

# The service listens on port 3000 inside the container
EXPOSE 3000

# Start the orderservice when the container runs
CMD ["/app/ordersystem"]
