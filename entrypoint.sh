#!/bin/sh

echo "Waiting for DB to be ready..."
sleep 5

echo "Running goose migrations..."
/go/bin/goose -dir /app/migrations/sql postgres "$DATABASE_DSN" up

echo "Starting main app..."
/app/main
