#!/bin/sh

# todo
# docker build
# docker run db
# docker run orderservice
#!/bin/sh
# ---------------------------------------------------------------------------
# Exercise 3 – Run script
# Builds and runs the PostgreSQL and OrderService containers on Linux/macOS.
# ---------------------------------------------------------------------------

set -e  # stop if any command fails

NETWORK="sbd3-net"
DB_CONTAINER="sbd3-postgres"
SERVICE_CONTAINER="orderservice"
DB_IMAGE="postgres:18"
SERVICE_IMAGE="orderservice:latest"
DB_VOLUME="sbd3-pgdata"

# 1) Ensure docker network exists
if ! docker network ls --format '{{.Name}}' | grep -q "^${NETWORK}$"; then
  echo ">> Creating network ${NETWORK}"
  docker network create "${NETWORK}"
fi

# 2) Start PostgreSQL
echo ">> Starting PostgreSQL container..."
docker run -d --name "${DB_CONTAINER}" \
  --network "${NETWORK}" \
  -e POSTGRES_DB=order \
  -e POSTGRES_USER=docker \
  -e POSTGRES_PASSWORD=docker123 \
  -v "${DB_VOLUME}:/var/lib/postgresql/18/docker" \
  "${DB_IMAGE}"

# 3) Build orderservice image
echo ">> Building OrderService image..."
docker build -t "${SERVICE_IMAGE}" .

# 4) Run orderservice (host 8080 → container 3000)
echo ">> Starting OrderService container..."
docker run -d --name "${SERVICE_CONTAINER}" \
  --network "${NETWORK}" \
  -p 8080:3000 \
  -e DB_HOST="${DB_CONTAINER}" \
  -e DB_PORT=5432 \
  -e DB_USER=docker \
  -e DB_PASSWORD=docker123 \
  -e DB_NAME=order \
  -e DB_SSLMODE=disable \
  -e HTTP_PORT=3000 \
  "${SERVICE_IMAGE}"

echo ""
echo "✅ Both containers are running."
echo "   Postgres  : ${DB_CONTAINER}"
echo "   Orderserve : http://localhost:8080/"
