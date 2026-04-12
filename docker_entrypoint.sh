#!/bin/sh
set -eu

/app/migrator -config /app/data/config.yml

exec /app/trash_project