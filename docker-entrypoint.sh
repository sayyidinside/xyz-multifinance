#!/bin/sh
set -e

# Wait for MySQL
while ! nc -z db 3306; do
  echo "Waiting for MySQL..."
  sleep 1
done

# Wait for Redis
while ! nc -z redis 6379; do
  echo "Waiting for Redis..."
  sleep 1
done

# Start the application
exec "$@"