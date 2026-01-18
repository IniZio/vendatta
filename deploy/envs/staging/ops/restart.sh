#!/bin/bash
set -e

cd "$(dirname "$0")/../.."
STAGING_DIR="${PWD}"
SERVER_ROOT="${STAGING_DIR%/*/*/*}"

echo "=== Restarting Nexus Coordination Server ==="
echo ""

PIDFILE="${STAGING_DIR}/.nexus/server.pid"

if [ -f "$PIDFILE" ]; then
    PID=$(cat "$PIDFILE")
    echo "Stopping server (PID: $PID)..."
    if kill -0 "$PID" 2>/dev/null; then
        kill "$PID"
        echo "Sent SIGTERM to server"
        sleep 2
        
        if kill -0 "$PID" 2>/dev/null; then
            echo "Force killing..."
            kill -9 "$PID"
        fi
    else
        echo "Process not running, removing stale PID file"
    fi
    rm -f "$PIDFILE"
else
    echo "No PID file found, attempting to find running server..."
    if pgrep -f "nexus coordination start" > /dev/null; then
        pkill -f "nexus coordination start" || true
        sleep 1
    fi
fi

echo "Waiting 2 seconds..."
sleep 2

echo "Starting server in background..."
export VENDETTA_COORD_CONFIG="${STAGING_DIR}/config/coordination.yaml"
export DB_PATH="${STAGING_DIR}/.nexus/nexus.db"

nohup bash -c "cd \"${SERVER_ROOT}\" && nexus coordination start" > "${STAGING_DIR}/.nexus/server.log" 2>&1 &
NEW_PID=$!
echo $NEW_PID > "$PIDFILE"

echo "Server started with PID: $NEW_PID"
echo "Log file: ${STAGING_DIR}/.nexus/server.log"
echo ""
echo "Monitoring server startup (15 seconds)..."
echo "---"

for i in {1..15}; do
    if curl -s http://localhost:3001/health > /dev/null 2>&1; then
        echo "✅ Server is healthy!"
        echo ""
        tail -20 "${STAGING_DIR}/.nexus/server.log"
        echo ""
        echo "=== Restart Complete ==="
        exit 0
    fi
    
    echo "Attempt $i/15... waiting for server..."
    sleep 1
done

echo "⚠️  Server not responding after 15 seconds"
echo ""
echo "Last 30 lines of log:"
tail -30 "${STAGING_DIR}/.nexus/server.log"
echo ""
echo "Monitor with: tail -f ${STAGING_DIR}/.nexus/server.log"
echo "Or check manually: curl http://localhost:3001/health"
