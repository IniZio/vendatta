#!/bin/bash
# M3 Transport Layer Verification Script
# Usage: ./scripts/verify-m3-transport.sh [--full]

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=============================================="
echo "  M3 Transport Layer Verification"
echo "=============================================="
echo ""

PASS=0
FAIL=0

# Function to print test result
print_result() {
    local test_name="$1"
    local result="$2"
    if [ "$result" -eq 0 ]; then
        echo -e "${GREEN}✅ PASS${NC}: $test_name"
        ((PASS++))
    else
        echo -e "${RED}❌ FAIL${NC}: $test_name"
        ((FAIL++))
    fi
}

# Check 1: Build succeeds
echo "1. Building project..."
if go build ./... 2>/dev/null; then
    print_result "Go build" 0
else
    print_result "Go build" 1
fi

# Check 2: Transport unit tests (short)
echo ""
echo "2. Running transport unit tests..."
if [ "$1" == "--full" ]; then
    go test ./pkg/transport/... -timeout 120s 2>/dev/null && print_result "Transport tests" 0 || print_result "Transport tests" 1
else
    # Quick smoke test
    go test ./pkg/transport/... -run "TestManager|TestSSHTransportCreation|TestHTTPTransportCreation" -timeout 30s 2>/dev/null && print_result "Transport smoke tests" 0 || print_result "Transport smoke tests" 1
fi

# Check 3: Coordination server unit tests
echo ""
echo "3. Running coordination unit tests..."
go test ./pkg/coordination/... -timeout 60s 2>/dev/null && print_result "Coordination tests" 0 || print_result "Coordination tests" 1

# Check 4: Agent unit tests
echo ""
echo "4. Running agent unit tests..."
go test ./pkg/agent/... -timeout 60s 2>/dev/null && print_result "Agent tests" 0 || print_result "Agent tests" 1

# Check 5: Provider tests (QEMU)
echo ""
echo "5. Running QEMU provider tests..."
go test ./pkg/provider/qemu/... -timeout 60s 2>/dev/null && print_result "QEMU provider tests" 0 || print_result "QEMU provider tests" 1

# Check 6: E2E test infrastructure
echo ""
echo "6. Checking E2E test infrastructure..."
if [ -f "e2e/transport_local_test.go" ] && [ -f "e2e/m3_verification_test.go" ]; then
    print_result "E2E test files exist" 0
else
    print_result "E2E test files exist" 1
fi

# Check 7: Verify transport interface
echo ""
echo "7. Verifying transport interface..."
if grep -q "type Transport interface" pkg/transport/interface.go 2>/dev/null; then
    print_result "Transport interface defined" 0
else
    print_result "Transport interface defined" 1
fi

# Check 8: Verify SSH transport
echo ""
echo "8. Verifying SSH transport..."
if grep -q "type SSHTransport struct" pkg/transport/ssh.go 2>/dev/null; then
    print_result "SSHTransport struct defined" 0
else
    print_result "SSHTransport struct defined" 1
fi

# Check 9: Verify HTTP transport
echo ""
echo "9. Verifying HTTP transport..."
if grep -q "type HTTPTransport struct" pkg/transport/http.go 2>/dev/null; then
    print_result "HTTPTransport struct defined" 0
else
    print_result "HTTPTransport struct defined" 1
fi

# Check 10: Verify connection pooling
echo ""
echo "10. Verifying connection pooling..."
if grep -q "type Pool struct" pkg/transport/pool.go 2>/dev/null; then
    print_result "Pool struct defined" 0
else
    print_result "Pool struct defined" 1
fi

# Summary
echo ""
echo "=============================================="
echo "  Summary"
echo "=============================================="
echo -e "Passed: ${GREEN}$PASS${NC}"
echo -e "Failed: ${RED}$FAIL${NC}"
echo ""

if [ $FAIL -eq 0 ]; then
    echo -e "${GREEN}All checks passed!${NC}"
    exit 0
else
    echo -e "${YELLOW}Some checks failed. Review output above.${NC}"
    exit 1
fi
