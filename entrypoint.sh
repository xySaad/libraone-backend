#!/bin/sh
set -e

echo "[entrypoint] applying database migrations..."
/app/migrate -source file:///app/migrations -database "sqlite3://db/database.db" up

echo "[entrypoint] starting libraone..."
exec /app/libraone