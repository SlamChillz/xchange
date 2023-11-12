#!/bin/sh

set -e

echo "run database migrations"
migrate -path /app/db/migrations -database "$DBURL" -verbose up

echo "start application"
exec "$@"
