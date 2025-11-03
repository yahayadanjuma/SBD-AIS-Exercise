# Dockerfile for Exercise 4 (folder: Exc_4/skeleton)
# Goal: build a small, secure image for the Go orderservice using a multi-stage build.
# Presenter note: Comments go ABOVE instructions (Dockerfile does not allow inline comments after commands).

# -------- Stage 1: Builder (compile the Go binary) --------
# Use a Go image that satisfies the go.mod toolchain requirement (>= 1.25.x).
FROM golang:1.25 AS builder

# Allow Go to auto-download the exact toolchain if "toolchain" is declared in go.mod.
ENV GOTOOLCHAIN=auto

# Create and move into a working directory inside the builder container.
WORKDIR /src

# Copy module files first to leverage Docker’s layer cache for dependency downloads.
# Presenter tip: If go.mod/go.sum didn’t change, this layer is re-used and speeds up rebuilds.
COPY go.mod go.sum ./

# Download all dependencies required by your module (cached when mod files unchanged).
RUN go mod download

# Copy the full source tree (the entire Exc_4/skeleton folder) AFTER deps are cached.
# Presenter note: This avoids re-downloading deps when only source files change.
COPY . .

# Build a static Linux binary for portability and a smaller final runtime image.
# - CGO_ENABLED=0 removes libc dependency (fully static)
# - GOOS/GOARCH targets Linux/amd64 (works on typical Docker hosts)
# - If your main.go lives in a subfolder (e.g., ./cmd/orderservice), change the build path accordingly.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /app/orderservice ./

# -------- Stage 2: Runtime (minimal, secure image) --------
# Tiny base image for running the compiled binary only (no compiler, no sources).
FROM alpine:3.20

# Set a clean working directory inside the runtime container.
WORKDIR /app

# Create a non-root user to run the service (security best practice).
RUN adduser -D appuser

# Copy only the compiled binary from the builder stage (keeps image small and safer).
COPY --from=builder /app/orderservice /app/orderservice

# Document the port the service listens on (Compose maps this to host: 3000:3000).
EXPOSE 3000

# Drop privileges and run as the non-root user created above.
USER appuser

# Start the service when the container launches.
# Presenter note: must match the path used in docker-compose.yml `command`.
ENTRYPOINT ["/app/orderservice"]
