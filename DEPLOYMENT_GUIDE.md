# SQLite Staging Deployment & Verification Guide

## Overview

This guide walks you through deploying the SQLite persistence changes to staging and verifying they work correctly.

**Key Features Being Tested:**
- ✅ Fork detection for private repositories
- ✅ Auto-forking to user account
- ✅ GitHub token persistence across restarts
- ✅ Fork information storage in SQLite
- ✅ Data persistence through server restart

---

## Phase 1: Deploy to Staging (Non-Blocking)

### Step 1: Verify Environment Setup

The staging environment is configured at: `deploy/envs/staging/`

**Confirm these are ready** (ask me if you need to check):
- GITHUB_APP_ID ✓
- GITHUB_APP_CLIENT_ID ✓
- GITHUB_APP_CLIENT_SECRET ✓
- GITHUB_APP_PRIVATE_KEY ✓
- GITHUB_APP_REDIRECT_URL ✓

### Step 2: Start Server with SQLite Enabled

```bash
cd deploy/envs/staging

# The binary is already built with SQLite support
# Enable SQLite by setting DB_PATH (restart.sh does this automatically)
./ops/restart.sh
```

**What this does:**
- Stops any running server gracefully
- Sets `DB_PATH=deploy/envs/staging/.nexus/nexus.db`
- Starts server in background (non-blocking)
- Monitors startup with 15-second health check
- Shows last 20 lines of log
- Returns immediately

**Expected Output:**
```
=== Restarting Nexus Coordination Server ===
...
✅ Server is healthy!
[Server startup logs...]
=== Restart Complete ===
```

### Step 3: Verify Server is Running

```bash
# In any terminal:
curl http://localhost:3001/health
```

Expected: `{"status":"ok"}`

### Step 4: Verify SQLite Database Created

```bash
# Check database file exists
ls -lh deploy/envs/staging/.nexus/nexus.db

# Check schema
sqlite3 deploy/envs/staging/.nexus/nexus.db ".tables"
```

Expected tables:
- `github_installations`
- `github_forks`
- `users`
- `workspaces`
- `services`
- `_schema_version`

---

## Phase 2: Step-by-Step Verification

### Prerequisites

You need:
1. **GitHub Account** with the app authorized
2. **Access to Private Repo** (e.g., oursky/epson-eshop)
3. **SSH Key** for workspace access

### Run Automated Verification Script

```bash
cd deploy/envs/staging
./ops/verify-sqlite.sh
```

This script will guide you through:
1. ✅ Verify database exists
2. ✅ Create a workspace
3. ✅ Check fork data in database
4. ✅ Stop workspace
5. ✅ Restart server
6. ✅ Verify data persists

---

## Phase 3: Manual E2E Testing

### Step 1: Create Workspace for Private Repo

```bash
cd deploy/envs/staging

# Register user (if not already done)
./ops/users.sh register IniZio 123456 "ssh-ed25519 AAAA..."

# Create workspace from private repo
./ops/workspaces.sh create IniZio test-sqlite-verification
```

**What happens:**
- System detects repo is private and not owned by user
- Automatically forks to user's account (`IniZio/epson-eshop`)
- Response includes:
  ```json
  {
    "workspace_id": "ws-1234567890",
    "fork_created": true,
    "fork_url": "https://github.com/IniZio/epson-eshop.git",
    "ssh_port": 2295,
    ...
  }
  ```

### Step 2: Verify Fork in Database

```bash
# Check fork was recorded
sqlite3 deploy/envs/staging/.nexus/nexus.db \
  "SELECT original_owner, original_repo, fork_owner FROM github_forks;"
```

Expected:
```
oursky|epson-eshop|IniZio
```

### Step 3: SSH into Workspace & Test Git

```bash
# Get SSH port from workspace status
./ops/workspaces.sh status ws-<ID>

# SSH into workspace
ssh -p 2295 dev@localhost

# Inside workspace:
cd /workspace

# Configure git
git config --global user.name "Test User"
git config --global user.email "test@example.com"

# Clone the fork (with token from environment)
git clone https://github.com/IniZio/epson-eshop.git epson-repo
cd epson-repo

# Make a test commit
echo "SQLite test: $(date)" >> SQLITE_TEST.txt
git add SQLITE_TEST.txt
git commit -m "test(sqlite): verify persistence works"
git push origin main

# Exit workspace
exit
```

### Step 4: Verify Commit on GitHub

Visit: `https://github.com/IniZio/epson-eshop/commits`

You should see your test commit from the workspace!

### Step 5: Stop Workspace

```bash
./ops/workspaces.sh stop ws-<ID>
```

### Step 6: Critical Test - Server Restart

```bash
# Restart server (data MUST persist)
./ops/restart.sh

# Wait for server to start
sleep 3
```

### Step 7: Verify Data Persists After Restart

```bash
# GitHub token should still be there
sqlite3 deploy/envs/staging/.nexus/nexus.db \
  "SELECT COUNT(*) FROM github_installations;"

# Fork info should still be there
sqlite3 deploy/envs/staging/.nexus/nexus.db \
  "SELECT fork_owner FROM github_forks LIMIT 1;"

# Check workspace still accessible
./ops/workspaces.sh list
./ops/workspaces.sh status ws-<ID>
```

**Expected Results:**
- ✅ GitHub installations: 1
- ✅ Fork owner: IniZio
- ✅ Workspace still listed
- ✅ Workspace status recoverable

---

## Troubleshooting

### Server Won't Start

```bash
# Check logs
tail -f deploy/envs/staging/.nexus/server.log

# Look for SQLite errors like:
# - "no such table: _schema_version" → migrations didn't run
# - "database is locked" → another process has it open
```

### Database Issues

```bash
# Check database integrity
sqlite3 deploy/envs/staging/.nexus/nexus.db "PRAGMA integrity_check;"

# Reset database (⚠️ deletes all data)
rm deploy/envs/staging/.nexus/nexus.db
./ops/restart.sh
```

### Fork Not Created

```bash
# Verify GitHub token is stored
sqlite3 deploy/envs/staging/.nexus/nexus.db \
  "SELECT user_id, github_username FROM github_installations;"

# Check if repo is actually private
curl -H "Authorization: token YOUR_TOKEN" \
  https://api.github.com/repos/oursky/epson-eshop | jq .private
```

### SSH Access Issues

```bash
# Get correct SSH port from workspace
./ops/workspaces.sh status ws-<ID> | grep ssh_port

# Test SSH
ssh -v -p 2295 dev@localhost

# If fails, check LXC is running
lxc list
```

---

## Verification Checklist

Print this and check off as you go:

```
[ ] Server starts with SQLite enabled
[ ] Database file created at .nexus/nexus.db
[ ] All tables exist (github_installations, github_forks, etc.)
[ ] Workspace created for private repo
[ ] Fork detected and auto-created
[ ] Fork info stored in database
[ ] SSH into workspace works
[ ] Git clone works with token
[ ] Git commit created
[ ] Git push succeeded
[ ] Commit visible on GitHub fork
[ ] Server restarted successfully
[ ] GitHub token persists after restart
[ ] Fork data persists after restart
[ ] Workspace still accessible after restart
```

---

## Cleanup

If you need to start fresh:

```bash
# Stop server
./ops/restart.sh  # Will stop first

# Delete database (starts fresh)
rm deploy/envs/staging/.nexus/nexus.db

# Delete workspaces (if needed)
./ops/workspaces.sh delete ws-<ID>

# Restart server
./ops/restart.sh
```

---

## Next Steps

After verification succeeds:

1. **Document Results** - Create a test report with screenshots
2. **Code Review** - Review changes:
   - `pkg/github/fork.go` - Fork API
   - `pkg/coordination/registry_sqlite.go` - Database persistence
   - `pkg/coordination/db.go` - Schema
3. **Merge to Main** - Create PR with test results
4. **Production Deploy** - Use same restart.sh in production

---

## Support

Issues? Questions? Create an issue or ask:
- Check logs: `tail -f deploy/envs/staging/.nexus/server.log`
- Manual health check: `curl http://localhost:3001/health`
- Database query: `sqlite3 deploy/envs/staging/.nexus/nexus.db`
