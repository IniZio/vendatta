#!/bin/bash
set -e

cd "$(dirname "$0")/../.."
STAGING_DIR="${PWD}"
SERVER_ROOT="${STAGING_DIR%/*/*/*}"

export VENDETTA_COORD_CONFIG="${STAGING_DIR}/config/coordination.yaml"

echo "=== Starting Nexus Coordination Server ==="
echo "Config: ${VENDETTA_COORD_CONFIG}"
echo "Database: ${STAGING_DIR}/.nexus/nexus.db"
echo "Server: http://localhost:3001"
echo ""
echo "To stop: Ctrl+C"
echo ""

cd "$SERVER_ROOT"
exec nexus coordination start
