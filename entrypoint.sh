#!/bin/sh

echo "⏳ Waiting for DB..."

export PGPASSWORD="$DB_PASSWORD"

# Wait for PostgreSQL to be ready
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" > /dev/null 2>&1; do
  echo "⌛ Still waiting for PostgreSQL at $DB_HOST:$DB_PORT..."
  sleep 2
done

echo "✅ PostgreSQL is ready."

# Wait for Keycloak OIDC config to be available
echo "⏳ Waiting for Keycloak to be ready..."
until curl -s "$KEYCLOAK_ISSUER/.well-known/openid-configuration" | grep -q '"issuer"'; do
  echo "⌛ Still waiting for Keycloak at $KEYCLOAK_ISSUER..."
  sleep 2
done
echo "✅ Keycloak is up and serving OIDC config."

echo "📂 Running migrations..."
/go/bin/goose -dir /app/migrations/sql postgres "$DATABASE_DSN" up

echo "🚀 Starting app..."
exec /app/main
