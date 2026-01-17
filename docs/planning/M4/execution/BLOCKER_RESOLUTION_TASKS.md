# Critical Blockers - Detailed Task Specifications

**Status**: Blocking Phase 1 start (Jan 27)  
**Priority**: CRITICAL  
**Assigned to**: Systems Engineer  
**Deadline**: Jan 22, 2026

---

## BLOCKER 1: Fix Failing E2E Tests

**Duration**: 8-12 hours  
**Deadline**: Jan 20, 2026

### Issue Summary

5 E2E test files are failing due to environment setup issues:
- `e2e/lifecycle_test.go` - Failing
- `e2e/m3_verification_test.go` - Failing
- `e2e/transport_test.go` - Failing
- `e2e/transport_local_test.go` - Failing
- `e2e/testenv.go` has test utilities but tests won't run

### Root Causes to Investigate

1. **Docker daemon not running or misconfigured**
   - Test: `docker ps` should return running containers
   - Fix: `sudo systemctl start docker` or equivalent

2. **SSH keys not available for tests**
   - Tests expect SSH keys at `~/.ssh/id_ed25519`
   - Fix: Generate test keys or configure CI to provide them

3. **LXC memory/resource limits causing failures**
   - Tests may fail if LXC containers can't allocate resources
   - Fix: Adjust LXC limits or skip LXC tests in CI

4. **Environment variables missing**
   - Tests may expect specific ENV vars for Docker/LXC
   - Fix: Set in CI configuration or local shell

5. **Test environment initialization incomplete**
   - `e2e/testenv.go` provides utilities but may need setup
   - Fix: Call setup functions before test execution

### Step-by-Step Resolution

**Step 1: Diagnose Current State** (2 hours)
```bash
# Run tests and capture actual errors
cd /home/newman/magic/nexus
make test-e2e 2>&1 | tee /tmp/e2e_errors.log

# Check Docker
docker ps
docker --version

# Check SSH keys
ls -la ~/.ssh/id_ed25519*

# Check LXC
lxc version
lxc list

# Check Go test environment
go test -v ./e2e/... -timeout 60s 2>&1 | head -100
```

**Step 2: Fix Docker Issues** (2 hours if needed)
```bash
# Ensure Docker daemon is running
sudo systemctl start docker
sudo usermod -aG docker $USER
newgrp docker

# Test Docker
docker run hello-world

# If still failing, check logs
journalctl -u docker -n 50
```

**Step 3: Fix SSH Key Issues** (1 hour if needed)
```bash
# Generate SSH keys if missing
[ -f ~/.ssh/id_ed25519 ] || ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519 -N ""

# Set proper permissions
chmod 600 ~/.ssh/id_ed25519
chmod 644 ~/.ssh/id_ed25519.pub

# Add to authorized_keys for local testing
cat ~/.ssh/id_ed25519.pub >> ~/.ssh/authorized_keys
chmod 600 ~/.ssh/authorized_keys
```

**Step 4: Fix LXC Issues** (2 hours if needed)
```bash
# Check LXC memory available
lxc info local: | grep memory

# Check LXC storage
lxc storage list

# Try launching test container
lxc launch ubuntu:22.04 test-container
lxc delete test-container -f

# If fails, check system resources
free -h
df -h
```

**Step 5: Run E2E Tests** (1 hour)
```bash
# Run with verbose output
go test -v ./e2e/... -timeout 120s 2>&1 | tee /tmp/e2e_results.log

# Expected output: "ok      github.com/nexus/nexus/e2e"
# All 5 test files should show PASS
```

### Success Criteria

- [ ] `make test-e2e` runs without errors
- [ ] All 5 E2E test files passing
- [ ] Tests complete within 120 seconds
- [ ] No environment setup required beyond standard tools
- [ ] CI/CD pipeline shows green for e2e tests

### Verification

Run this command and verify output:
```bash
cd /home/newman/magic/nexus
go test -v ./e2e/... -timeout 120s | grep -E "^(ok|FAIL|--- PASS|--- FAIL)"
```

Expected output:
```
--- PASS: TestLifecycle (30s)
--- PASS: TestM3Verification (25s)
--- PASS: TestTransport (20s)
--- PASS: TestTransportLocal (15s)
ok  github.com/nexus/nexus/e2e  120s
```

### Documentation

When complete, create `docs/environment-setup.md` with:
- Docker setup instructions
- SSH key requirements
- LXC resource requirements
- How to run tests locally
- Troubleshooting common issues

---

## BLOCKER 2: Implement Makefile Test Targets

**Duration**: 4-6 hours  
**Deadline**: Jan 22, 2026

### Issue Summary

Makefile declares test targets but implementations are missing:

Current state (non-functional):
```makefile
test:
	@echo "Running tests..."
test-unit:
	@echo "Running unit tests..."
test-integration:
	@echo "Running integration tests..."
test-e2e:
	@echo "Running e2e tests..."
```

### Required Implementation

**Step 1: Open Makefile**
```bash
cd /home/newman/magic/nexus
cat Makefile | head -50
```

**Step 2: Implement Test Targets**

Replace the placeholder targets with actual implementations:

```makefile
# Test targets
.PHONY: test test-unit test-integration test-e2e test-coverage test-all

# Run all tests
test:
	@echo "Running all tests..."
	@go test -v -race -coverprofile=coverage.out ./... -timeout 120s

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	@go test -v -race -short -coverprofile=coverage.unit.out ./... -timeout 60s

# Run integration tests only (exclude e2e)
test-integration:
	@echo "Running integration tests..."
	@go test -v -race -coverprofile=coverage.integration.out \
		-run "Integration" ./... -timeout 120s

# Run e2e tests only
test-e2e:
	@echo "Running e2e tests..."
	@go test -v ./e2e/... -timeout 120s

# Generate coverage report
test-coverage:
	@echo "Generating coverage report..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run all tests with coverage
test-all: test test-coverage
	@echo "All tests completed. Coverage: coverage.html"
```

**Step 3: Add Helper Targets**

Add these supporting targets:

```makefile
# Clean test artifacts
test-clean:
	@rm -f coverage*.out coverage.html

# Run tests with race detector
test-race:
	@go test -v -race ./... -timeout 120s

# Run specific test
test-run:
	@go test -v -run $(TEST) ./... -timeout 120s
```

**Step 4: Test the Targets**

```bash
# Test each target
make test-unit          # Should run unit tests
make test-integration   # Should run integration tests
make test-e2e          # Should run e2e tests
make test-coverage     # Should generate coverage.html
make test              # Should run all tests
```

### Success Criteria

- [ ] `make test` runs all tests successfully
- [ ] `make test-unit` runs unit tests only
- [ ] `make test-integration` runs integration tests
- [ ] `make test-e2e` runs E2E tests
- [ ] `make test-coverage` generates coverage report
- [ ] All targets have proper error handling
- [ ] Coverage report viewable in browser

### Verification

```bash
# Each should complete without errors
make test-unit
echo "Exit code: $?"  # Should be 0

make test-integration
echo "Exit code: $?"  # Should be 0

make test-e2e
echo "Exit code: $?"  # Should be 0

make test
echo "Exit code: $?"  # Should be 0
```

### Update CI/CD Pipeline

Update `.github/workflows/ci.yaml` to use the new targets:

```yaml
name: CI

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      
      - name: Run unit tests
        run: make test-unit
      
      - name: Run integration tests
        run: make test-integration
      
      - name: Run e2e tests
        run: make test-e2e
      
      - name: Generate coverage
        run: make test-coverage
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out
```

### Documentation

Add to README.md:

```markdown
## Testing

Run tests locally:

```bash
make test              # Run all tests
make test-unit        # Unit tests only
make test-integration # Integration tests only
make test-e2e         # E2E tests only
make test-coverage    # Generate coverage report
```

Tests require:
- Go 1.24+
- Docker (for some tests)
- SSH keys (~/.ssh/id_ed25519)
- LXC (for container tests)
```

---

## VERIFICATION CHECKLIST (Jan 22)

Before Phase 1 starts, verify:

- [ ] BLOCKER 1: All E2E tests passing
  ```bash
  go test -v ./e2e/... -timeout 120s
  # Should show: ok  github.com/nexus/nexus/e2e
  ```

- [ ] BLOCKER 2: All Makefile targets working
  ```bash
  make test
  # Should complete without errors
  ```

- [ ] CI/CD green
  ```bash
  # Check GitHub Actions: all workflows passing
  ```

- [ ] Documentation updated
  ```bash
  # docs/environment-setup.md created
  # README.md updated with test instructions
  ```

---

## Handoff to Phase 1

Once both blockers are fixed:

1. Push changes to main branch
2. Verify all CI checks passing
3. Notify Sisyphus: "Blockers resolved, ready for Phase 1"
4. Schedule Phase 1 kickoff (Jan 27)

---

**Owner**: Systems Engineer  
**Status**: Not Started  
**Due**: Jan 22, 2026  
**Estimated Effort**: 12-18 hours
