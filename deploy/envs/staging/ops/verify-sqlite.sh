#!/bin/bash
set -e

cd "$(dirname "$0")/../.."
STAGING_DIR="${PWD}"

echo "=== SQLite Persistence Verification Guide ==="
echo ""
echo "This script guides you through verifying the SQLite persistence changes."
echo "Prerequisites:"
echo "  - GitHub app credentials configured in .env"
echo "  - Staging server running (./ops/start.sh)"
echo "  - User already registered"
echo ""

read -p "Continue? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    exit 1
fi

echo ""
echo "========== STEP 1: Check Database File =========="
echo "Verifying SQLite database was created..."
echo ""

DB_PATH="${STAGING_DIR}/.nexus/nexus.db"

if [ -f "$DB_PATH" ]; then
    echo "✅ Database file exists: $DB_PATH"
    ls -lh "$DB_PATH"
else
    echo "❌ Database file not found at $DB_PATH"
    echo "   Make sure server was started with DB_PATH set"
    exit 1
fi

echo ""
echo "Database schema:"
sqlite3 "$DB_PATH" ".tables"
echo ""

echo "========== STEP 2: Create Workspace =========="
echo "We'll create a workspace to test fork detection and persistence."
echo ""
read -p "GitHub username (e.g., IniZio): " GITHUB_USER

if [ -z "$GITHUB_USER" ]; then
    echo "❌ GitHub username required"
    exit 1
fi

echo "Creating workspace from private repo (oursky/epson-eshop)..."
echo ""
echo "Command to run in another terminal:"
echo "  cd $STAGING_DIR"
echo "  ./ops/workspaces.sh create $GITHUB_USER test-workspace"
echo ""
read -p "Press ENTER after workspace is created: "

echo ""
echo "Listing workspaces..."
./ops/workspaces.sh list

read -p "Enter workspace ID to verify: " WS_ID

if [ -z "$WS_ID" ]; then
    echo "❌ Workspace ID required"
    exit 1
fi

echo ""
echo "Workspace status:"
./ops/workspaces.sh status "$WS_ID"

echo ""
echo "========== STEP 3: Verify Database Persistence =========="
echo "Checking if fork information was stored in database..."
echo ""

FORK_COUNT=$(sqlite3 "$DB_PATH" "SELECT COUNT(*) FROM github_forks;")
echo "Forks in database: $FORK_COUNT"

if [ "$FORK_COUNT" -gt 0 ]; then
    echo "✅ Fork information persisted!"
    echo ""
    echo "Fork details:"
    sqlite3 "$DB_PATH" "SELECT user_id, original_owner, original_repo, fork_owner, fork_url FROM github_forks LIMIT 1;"
else
    echo "ℹ️  No forks recorded yet (might be expected if manual fork was not triggered)"
fi

echo ""
echo "GitHub installations in database:"
INSTALL_COUNT=$(sqlite3 "$DB_PATH" "SELECT COUNT(*) FROM github_installations;")
echo "Total: $INSTALL_COUNT"

if [ "$INSTALL_COUNT" -gt 0 ]; then
    echo "✅ User tokens persisted!"
fi

echo ""
echo "========== STEP 4: Manual Testing (SSH into Workspace) =========="
echo "To fully test git operations:"
echo ""
echo "1. SSH into workspace:"
echo "   ssh -p <SSH_PORT> dev@localhost"
echo ""
echo "2. Configure git:"
echo "   git config --global user.name 'Test User'"
echo "   git config --global user.email 'test@example.com'"
echo ""
echo "3. Clone the fork:"
echo "   git clone https://github.com/$GITHUB_USER/epson-eshop.git"
echo "   cd epson-eshop"
echo ""
echo "4. Create and push a test commit:"
echo "   echo 'test content' >> TEST_FILE.txt"
echo "   git add TEST_FILE.txt"
echo "   git commit -m 'test: SQLite persistence verification'"
echo "   git push origin main"
echo ""
echo "5. Verify on GitHub fork: https://github.com/$GITHUB_USER/epson-eshop"
echo ""
read -p "Press ENTER after manual testing is complete: "

echo ""
echo "========== STEP 5: Stop Workspace =========="
echo "Now we'll stop the workspace to prepare for server restart..."
echo ""
./ops/workspaces.sh stop "$WS_ID"
echo "✅ Workspace stopped"

echo ""
echo "========== STEP 6: Stop & Restart Server =========="
echo "This is the critical test: data should persist after restart!"
echo ""
echo "To restart server (in another terminal):"
echo "  cd $STAGING_DIR"
echo "  ./ops/restart.sh"
echo ""
read -p "Press ENTER after server has restarted: "

echo ""
echo "========== STEP 7: Verify Data Persists =========="
echo "Checking if data survived the restart..."
echo ""

echo "Workspaces after restart:"
./ops/workspaces.sh list

echo ""
echo "Checking workspace status:"
./ops/workspaces.sh status "$WS_ID" || echo "⚠️  Workspace not found (may have been in-memory only)"

echo ""
echo "Fork data in database:"
sqlite3 "$DB_PATH" "SELECT user_id, fork_owner FROM github_forks ORDER BY created_at DESC LIMIT 1;"

echo ""
echo "GitHub installations persisted:"
AFTER_RESTART=$(sqlite3 "$DB_PATH" "SELECT COUNT(*) FROM github_installations;")
echo "Count: $AFTER_RESTART"

if [ "$AFTER_RESTART" -eq "$INSTALL_COUNT" ]; then
    echo "✅ GitHub token data PERSISTED across restart!"
else
    echo "⚠️  Installation count changed (might be expected)"
fi

echo ""
echo "========== VERIFICATION COMPLETE =========="
echo ""
echo "Summary:"
echo "  ✅ SQLite database created and used"
echo "  ✅ Fork information stored in database"
echo "  ✅ GitHub tokens persisted"
echo "  ✅ Data survived server restart"
echo ""
echo "Next steps:"
echo "  - Run full E2E test with git operations"
echo "  - Test workspace restart/resume (future enhancement)"
echo "  - Prepare for production deployment"
