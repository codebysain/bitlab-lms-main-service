#!/bin/sh
echo "⏳ Waiting for DB..."
sleep 10

echo "📂 Running migrations..."
/go/bin/goose -dir /app/migrations/sql postgres "$DATABASE_DSN" up

echo "🚀 Starting app..."
exec /app/main
