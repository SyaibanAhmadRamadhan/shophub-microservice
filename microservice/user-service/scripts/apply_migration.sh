#!/usr/bin/env bash
# This script runs database migrations using golang-migrate CLI.
# Works in local docker-compose environment.

set -e

# Default values
DATABASE_DSN="postgres://adminsi:adminsi_pw@localhost:5433/adminsi?sslmode=disable"
MIGRATE_VERSIONS=""  # Optional specific version
MIGRATIONS_DIR="$(pwd)/migrations"

# Parse CLI options
OPTSTRING=":d:v:"

while getopts ${OPTSTRING} opt; do
  case ${opt} in
    d)
      DATABASE_DSN=${OPTARG}
      ;;
    v)
      MIGRATE_VERSIONS=${OPTARG}
      ;;
    :)
      echo "❌ Option -${OPTARG} requires an argument." >&2
      exit 1
      ;;
    ?)
      echo "❌ Invalid option: -${OPTARG}." >&2
      exit 1
      ;;
  esac
done

# Check for required tools
command -v migrate >/dev/null 2>&1 || {
  echo "❌ 'migrate' CLI not found. Install via 'brew install golang-migrate' or from https://github.com/golang-migrate/migrate" >&2
  exit 1
}
command -v pgsanity >/dev/null 2>&1 || {
  echo "❌ 'pgsanity' not found. Install via pip or package manager." >&2
  exit 1
}

echo "🔍 Checking SQL files with pgsanity..."
find ./migrations -name '*.sql' ! -name 'seed_*' | xargs pgsanity

if [ -n "$MIGRATE_VERSIONS" ]; then
  echo "⏫ Migrating up to version $MIGRATE_VERSIONS..."
  migrate -database "$DATABASE_DSN" -path "$MIGRATIONS_DIR" goto "$MIGRATE_VERSIONS"
else
  echo "⏫ Running full migration (up)..."
  migrate -database "$DATABASE_DSN" -path "$MIGRATIONS_DIR" up
fi

echo "✅ Migration completed."
