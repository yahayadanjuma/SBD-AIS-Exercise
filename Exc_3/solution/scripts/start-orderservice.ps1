<# 
  start-orderservice.ps1
  Build and run the Order Service container, connected to Postgres.
  Works on Windows PowerShell 5.1 and PowerShell 7+.

  Defaults (edit if yours differ):
    - Docker network: sbd3-net
    - Postgres container: sbd3-postgres
    - Service image: orderservice
    - Service container: orderservice
    - App port (host & container): 3000
    - Optional env file beside this script: debug.env
#>

param(
  [int]    $AppPort       = 3000,
  [string] $ImageName     = "orderservice",
  [string] $ContainerName = "orderservice",
  [string] $Network       = "sbd3-net",
  [string] $PgContainer   = "sbd3-postgres",
  [string] $EnvFile       = "debug.env"
)

$ErrorActionPreference = "Stop"

function Exists($cmd) {
  $null -ne (Get-Command $cmd -ErrorAction SilentlyContinue)
}

function Ensure-Network([string]$name) {
  if (-not (docker network ls --format '{{.Name}}' | Select-String -SimpleMatch $name)) {
    Write-Host ">> Creating docker network '$name'..."
    & docker network create $name | Out-Null
    if ($LASTEXITCODE -ne 0) { throw "Failed to create network '$name'." }
  } else {
    Write-Host ">> Network '$name' already exists."
  }
}

function Ensure-Postgres([string]$name) {
  $running = & docker ps --filter "name=$name" --filter "status=running" --format '{{.Names}}'
  if (-not $running) {
    throw "PostgreSQL container '$name' is not running. Start DB first."
  } else {
    Write-Host ">> PostgreSQL '$name' is running."
  }
}

function Load-Env([string]$path, [hashtable]$defaults) {
  if (Test-Path $path) {
    Write-Host ">> Loading env from $path"
    $ht = @{}
    Get-Content $path | ForEach-Object {
      if (-not $_) { return }
      $line = $_.Trim()
      if ($line.StartsWith("#")) { return }
      $parts = $line -split "=", 2
      if ($parts.Count -eq 2) {
        $ht[$parts[0].Trim()] = $parts[1].Trim()
      }
    }
    foreach ($k in $defaults.Keys) {
      if (-not $ht.ContainsKey($k)) { $ht[$k] = $defaults[$k] }
    }
    return $ht
  }
  Write-Host ">> No env file found; using defaults."
  return $defaults
}

function Build-Image([string]$image) {
  Write-Host ">> docker build -t $image ."
  & docker build -t $image .
  if ($LASTEXITCODE -ne 0) { throw "Build failed." }
}

function Stop-Old([string]$name) {
  $exists = & docker ps -a --filter "name=$name" --format '{{.ID}}'
  if ($exists) {
    Write-Host ">> Removing existing container '$name'..."
    & docker rm -f $name | Out-Null
    if ($LASTEXITCODE -ne 0) { throw "Failed to remove existing container '$name'." }
  }
}

function Run-Service(
  [string]$image,
  [string]$name,
  [string]$net,
  [int]$portHost,
  [int]$portCont,
  [hashtable]$envs
) {
  # Build a safe argv list for: docker run ...
  $args = @(
    "run", "-d",
    "--name", $name,
    "--restart", "unless-stopped",
    "--network", $net,
    "-p", "$($portHost)`:$($portCont)"  # escape ':' to avoid $var:scope parsing in PS
  )

  foreach ($k in $envs.Keys) {
    $args += "-e"
    $args += "$k=$($envs[$k])"
  }

  $args += $image

  Write-Host ">> docker $($args -join ' ')"
  & docker @args | Out-Null
  if ($LASTEXITCODE -ne 0) { throw "docker run failed." }
}

# --------------------- MAIN ---------------------
if (-not (Exists "docker")) { throw "Docker CLI not found in PATH." }

Ensure-Network -name $Network
Ensure-Postgres -name $PgContainer

$defaults = @{
  "DB_HOST"     = $PgContainer
  "DB_PORT"     = "5432"
  "DB_USER"     = "postgres"
  "DB_PASSWORD" = "postgres"   # <-- set your real password
  "DB_NAME"     = "orders"
  "DB_SSLMODE"  = "disable"
  "HTTP_PORT"   = "$AppPort"
}

$envs = Load-Env -path $EnvFile -defaults $defaults

Build-Image -image $ImageName
Stop-Old -name $ContainerName
Run-Service -image $ImageName -name $ContainerName -net $Network -portHost $AppPort -portCont 3000 -envs $envs


Write-Host ""
Write-Host ">> Last 50 logs:"
& docker logs --tail 50 $ContainerName
Write-Host ""
Write-Host "✅ Done. Open http://localhost:$AppPort/"
