#!/bin/bash

echo "⏳ Waiting for MinIO to be ready..."
/wait-for-it.sh minio:9000 --timeout=30 --strict -- echo "✅ MinIO is up!"

echo "⏳ Waiting for PostgreSQL to be ready..."
/wait-for-it.sh app_db:5432 --timeout=30 --strict -- echo "✅ PostgreSQL is up!"

echo "🚀 Running goose migrations..."
/go/bin/goose -dir /app/migrations/sql postgres "$DATABASE_DSN" up

echo "🎬 Starting main app..."
exec /app/main
