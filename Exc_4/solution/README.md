Software Architecture for Big Data — Exercise 4 (Solution)
What this runs

orderservice: Go HTTP service (port 3000)

postgres: PostgreSQL 16 (port 5432)

Shared user-defined network, DB healthcheck, and a named volume for persistence.

Files

docker-compose.yml — Orchestrates both containers (network, ports, env, healthcheck)

Dockerfile — Multi-stage build (Go 1.25 builder → minimal Alpine runtime)

.env — DB credentials and config (dev defaults)

Prerequisites

Docker Desktop with Compose

Ports 3000 and 5432 free on the host

How to run
cd Exc_4\solution
docker compose up --build -d

Verify
docker compose ps
docker compose logs -f postgres          # wait for “database system is ready…”
docker compose logs -f orderservice      # should show “Order System is up and running”
curl.exe http://localhost:3000/          # or your app’s health route if available

DB sanity checks
docker compose exec postgres psql -U docker -d order -c "SELECT 1 as ok;"
docker compose exec postgres psql -U docker -d order -c "SELECT * FROM orders;"

Persistence check
docker compose down
docker compose up -d
docker compose exec postgres psql -U docker -d order -c "SELECT * FROM orders;"

Notes / gotchas (what we fixed)

Network/Container name clashes: use a unique network (ordersystem-ex4-net) and avoid hard-coded container_name for Postgres.

Go toolchain: golang:1.25 and GOTOOLCHAIN=auto satisfy go.mod.

App migration error (drink_id NOT NULL): either truncate demo rows or backfill drink_id before enforcing NOT NULL.

Bind address: the service must listen on :3000 (not localhost:3000) inside the container.