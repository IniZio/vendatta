#!/bin/bash
set -e

cd "$(dirname "$0")/../.."
STAGING_DIR="${PWD}"
SERVER_ROOT="${STAGING_DIR%/*/*/*}"

export VENDETTA_COORD_CONFIG="${STAGING_DIR}/config/coordination.yaml"

cd "$SERVER_ROOT"
exec nexus coordination restart
