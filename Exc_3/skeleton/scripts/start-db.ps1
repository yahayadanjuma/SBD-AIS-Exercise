# ===============================================================
# start-db.ps1
# Purpose: Starts PostgreSQL Docker container for SBD Exercise 3
# Compatible with Windows PowerShell
# ===============================================================

$ErrorActionPreference = 'Stop'

Write-Host "=== Starting PostgreSQL for SBD3 ===`n"

# 1) Ensure Docker network exists
if (-not (docker network ls --format '{{.Name}}' | Select-String -Quiet '^sbd3-net$')) {
    docker network create sbd3-net | Out-Null
    Write-Host "✅ Created Docker network: sbd3-net"
} else {
    Write-Host "ℹ️  Docker network 'sbd3-net' already exists"
}

# 2) Ensure persistent volume exists
if (-not (docker volume ls --format '{{.Name}}' | Select-String -Quiet '^sbd3-pgdata$')) {
    docker volume create sbd3-pgdata | Out-Null
    Write-Host "✅ Created Docker volume: sbd3-pgdata"
} else {
    Write-Host "ℹ️  Docker volume 'sbd3-pgdata' already exists"
}

# 3) Remove old Postgres container if it exists
if (docker ps -a --format '{{.Names}}' | Select-String -Quiet '^sbd3-postgres$') {
    Write-Host "🧹 Removing old 'sbd3-postgres' container..."
    docker rm -f sbd3-postgres | Out-Null
}

# 4) Start PostgreSQL container
Write-Host "🚀 Starting new PostgreSQL container..."
docker run -d --name sbd3-postgres `
  --network sbd3-net `
  -p 5432:5432 `
  -e POSTGRES_DB=order `
  -e POSTGRES_USER=docker `
  -e POSTGRES_PASSWORD=docker123 `
  -v sbd3-pgdata:/var/lib/postgresql/data `
  postgres:18 | Out-Null

Start-Sleep -Seconds 3
Write-Host "⏳ Waiting for Postgres to accept connections..."

# 5) Wait until Postgres is ready
$maxTries = 60
$ready = $false
for ($i = 1; $i -le $maxTries; $i++) {
    try {
        $result = docker exec sbd3-postgres bash -lc "PGPASSWORD=docker123 pg_isready -h localhost -p 5432 -U docker" 2>$null
        if ($result -match "accepting connections") {
            $ready = $true
            break
        }
    } catch {
        # just wait
    }
    Start-Sleep -Seconds 1
}

if ($ready) {
    Write-Host "`n PostgreSQL is ready and accepting connections on localhost:5432"
    Write-Host "   Container: sbd3-postgres"
} else {
    Write-Host "`n PostgreSQL did not start in time. Check logs:"
    docker logs sbd3-postgres
    exit 1
}

Write-Host "`nAll good! 🎉"
