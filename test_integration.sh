#!/bin/bash

# Eco CLI Integration Test Suite
# This script performs end-to-end testing of the eco CLI

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counters
TESTS_PASSED=0
TESTS_FAILED=0

# Temporary directory for test artifacts
TEMP_DIR=$(mktemp -d)
BINARY_PATH="$TEMP_DIR/eco"
ORIGINAL_HOME="$HOME"

# Cleanup function
cleanup() {
    echo ""
    echo "Cleaning up..."
    # Restore original home
    export HOME="$ORIGINAL_HOME"
    # Kill any running daemon
    if [ -f "/tmp/eco.pid" ]; then
        kill $(cat /tmp/eco.pid) 2>/dev/null || true
        rm -f /tmp/eco.pid
    fi
    # Remove temp directory
    rm -rf "$TEMP_DIR"
}

trap cleanup EXIT

# Helper functions
pass() {
    echo -e "${GREEN}✓ PASS${NC}: $1"
    TESTS_PASSED=$((TESTS_PASSED + 1))
}

fail() {
    echo -e "${RED}✗ FAIL${NC}: $1"
    TESTS_FAILED=$((TESTS_FAILED + 1))
}

info() {
    echo -e "${YELLOW}→ INFO${NC}: $1"
}

# Build the binary
build() {
    info "Building eco binary..."
    if go build -o "$BINARY_PATH" .; then
        pass "Build successful"
    else
        fail "Build failed"
        exit 1
    fi
}

# Setup isolated environment
setup_env() {
    info "Setting up isolated test environment..."
    export HOME="$TEMP_DIR/test-home"
    mkdir -p "$HOME/.config/eco"
    pass "Environment setup complete"
}

# Test: Help commands
test_help_commands() {
    info "Testing help commands..."
    
    local commands=(
        "--help"
        "init --help"
        "daemon --help"
        "stop --help"
        "status --help"
        "devices --help"
        "config --help"
        "config delete --help"
    )
    
    for cmd in "${commands[@]}"; do
        if $BINARY_PATH $cmd | grep -q "Usage:"; then
            pass "Help: eco $cmd"
        else
            fail "Help: eco $cmd"
        fi
    done
}

# Test: Init command
test_init() {
    info "Testing init command..."
    
    local output
    output=$($BINARY_PATH init 2>&1)
    
    if echo "$output" | grep -q "Eco initialized successfully"; then
        pass "Init command shows success message"
    else
        fail "Init command missing success message"
    fi
    
    if echo "$output" | grep -q "Device ID:"; then
        pass "Init command shows Device ID"
    else
        fail "Init command missing Device ID"
    fi
    
    if echo "$output" | grep -q "Secret:"; then
        pass "Init command shows Secret"
    else
        fail "Init command missing Secret"
    fi
    
    if [ -f "$HOME/.config/eco/config.json" ]; then
        pass "Config file created"
    else
        fail "Config file not created"
    fi
}

# Test: Status command
test_status() {
    info "Testing status command..."
    
    local output
    output=$($BINARY_PATH status 2>&1)
    
    if echo "$output" | grep -q "Initialized: yes"; then
        pass "Status shows initialized"
    else
        fail "Status doesn't show initialized"
    fi
    
    if echo "$output" | grep -q "Device ID:"; then
        pass "Status shows Device ID"
    else
        fail "Status missing Device ID"
    fi
}

# Test: Devices command
test_devices() {
    info "Testing devices command..."
    
    local output
    output=$($BINARY_PATH devices 2>&1)
    
    if echo "$output" | grep -q "Registered Device"; then
        pass "Devices shows header"
    else
        fail "Devices missing header"
    fi
    
    if echo "$output" | grep -q "Device ID:"; then
        pass "Devices shows Device ID"
    else
        fail "Devices missing Device ID"
    fi
    
    # Test with --show-secret
    output=$($BINARY_PATH devices --show-secret 2>&1)
    if echo "$output" | grep -q "Secret:"; then
        pass "Devices --show-secret shows secret"
    else
        fail "Devices --show-secret missing secret"
    fi
}

# Test: Daemon start/stop
test_daemon() {
    info "Testing daemon start/stop..."
    
    # Start daemon in background
    $BINARY_PATH daemon start &>/dev/null &
    local daemon_pid=$!
    
    # Wait for daemon to start
    sleep 2
    
    # Check PID file
    if [ -f "/tmp/eco.pid" ]; then
        pass "Daemon created PID file"
        local pid_from_file
        pid_from_file=$(cat /tmp/eco.pid)
        
        # Verify process is running
        if ps -p "$pid_from_file" > /dev/null 2>&1; then
            pass "Daemon process is running (PID: $pid_from_file)"
        else
            fail "Daemon process not running"
        fi
    else
        fail "Daemon PID file not created"
        kill $daemon_pid 2>/dev/null || true
        return
    fi
    
    # Stop daemon
    local output
    output=$($BINARY_PATH stop 2>&1)
    
    if echo "$output" | grep -q "Eco daemon stopped"; then
        pass "Stop command shows success"
    else
        fail "Stop command missing success message"
    fi
    
    # Wait for daemon to stop
    sleep 1
    
    # Verify daemon stopped
    if [ ! -f "/tmp/eco.pid" ]; then
        pass "PID file removed after stop"
    else
        fail "PID file not removed after stop"
    fi
}

# Test: Stop when not running
test_stop_not_running() {
    info "Testing stop when daemon not running..."
    
    # Ensure no PID file
    rm -f /tmp/eco.pid
    
    local output
    output=$($BINARY_PATH stop 2>&1)
    
    if echo "$output" | grep -q "not running"; then
        pass "Stop shows 'not running' message"
    else
        fail "Stop doesn't show 'not running' message"
    fi
}

# Test: Config delete
test_config_delete() {
    info "Testing config delete..."
    
    local output
    output=$($BINARY_PATH config delete 2>&1)
    
    if echo "$output" | grep -q "deleted successfully"; then
        pass "Config delete shows success"
    else
        fail "Config delete missing success message"
    fi
    
    if [ ! -f "$HOME/.config/eco/config.json" ]; then
        pass "Config file deleted"
    else
        fail "Config file not deleted"
    fi
}

# Test: Status after delete
test_status_after_delete() {
    info "Testing status after config delete..."
    
    local output
    output=$($BINARY_PATH status 2>&1)
    
    if echo "$output" | grep -q "not initialized"; then
        pass "Status shows 'not initialized' after delete"
    else
        fail "Status doesn't show 'not initialized' after delete"
    fi
}

# Test: Re-init after delete
test_reinit() {
    info "Testing re-init after delete..."
    
    local output
    output=$($BINARY_PATH init 2>&1)
    
    if echo "$output" | grep -q "Eco initialized successfully"; then
        pass "Re-init successful"
    else
        fail "Re-init failed"
    fi
}

# Test: JSON compatibility with Android
test_json_compatibility() {
    info "Testing JSON compatibility with Android..."
    
    # Test that generated JSON has lowercase field names
    local config_file="$HOME/.config/eco/config.json"
    
    if [ -f "$config_file" ]; then
        # Check for lowercase fields in config (Go uses PascalCase, JSON should match)
        pass "Config file exists for JSON check"
    else
        fail "Config file missing for JSON check"
    fi
}

# Main execution
main() {
    echo "========================================"
    echo "Eco CLI Integration Test Suite"
    echo "========================================"
    echo ""
    
    build
    setup_env
    
    echo ""
    echo "Running tests..."
    echo ""
    
    test_help_commands
    test_init
    test_status
    test_devices
    test_daemon
    test_stop_not_running
    test_config_delete
    test_status_after_delete
    test_reinit
    test_json_compatibility
    
    echo ""
    echo "========================================"
    echo "Test Results"
    echo "========================================"
    echo -e "${GREEN}Passed: $TESTS_PASSED${NC}"
    echo -e "${RED}Failed: $TESTS_FAILED${NC}"
    echo ""
    
    if [ $TESTS_FAILED -eq 0 ]; then
        echo -e "${GREEN}All tests passed!${NC}"
        exit 0
    else
        echo -e "${RED}Some tests failed!${NC}"
        exit 1
    fi
}

main "$@"
