#!/bin/bash

# Eco CLI Test Runner
# Runs all tests including unit tests and integration tests

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}    Eco CLI Test Runner${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Function to run tests
run_test() {
    local name="$1"
    local cmd="$2"
    
    echo -e "${YELLOW}Running: $name${NC}"
    echo "Command: $cmd"
    echo ""
    
    if eval "$cmd"; then
        echo ""
        echo -e "${GREEN}✓ $name passed${NC}"
    else
        echo ""
        echo -e "${RED}✗ $name failed${NC}"
        return 1
    fi
    echo ""
}

# Parse arguments
RUN_UNIT=true
RUN_INTEGRATION=false
RUN_BASH=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --unit|-u)
            RUN_UNIT=true
            RUN_INTEGRATION=false
            RUN_BASH=false
            shift
            ;;
        --integration|-i)
            RUN_UNIT=false
            RUN_INTEGRATION=true
            RUN_BASH=false
            shift
            ;;
        --bash|-b)
            RUN_UNIT=false
            RUN_INTEGRATION=false
            RUN_BASH=true
            shift
            ;;
        --all|-a)
            RUN_UNIT=true
            RUN_INTEGRATION=true
            RUN_BASH=true
            shift
            ;;
        --help|-h)
            echo "Usage: $0 [options]"
            echo ""
            echo "Options:"
            echo "  --unit, -u         Run only unit tests"
            echo "  --integration, -i  Run only integration tests (Go)"
            echo "  --bash, -b         Run only bash integration tests"
            echo "  --all, -a          Run all tests (default: unit only)"
            echo "  --help, -h         Show this help"
            echo ""
            echo "Examples:"
            echo "  $0                 # Run unit tests only"
            echo "  $0 --all           # Run all tests"
            echo "  $0 --integration   # Run Go integration tests"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Track overall success
ALL_PASSED=true

# Run unit tests
if [ "$RUN_UNIT" = true ]; then
    echo -e "${BLUE}----------------------------------------${NC}"
    echo -e "${BLUE}Unit Tests${NC}"
    echo -e "${BLUE}----------------------------------------${NC}"
    echo ""
    
    if ! run_test "Protocol Tests" "go test ./internal/protocol/... -v"; then
        ALL_PASSED=false
    fi
    
    if ! run_test "Config Tests" "go test ./internal/config/... -v"; then
        ALL_PASSED=false
    fi
fi

# Run Go integration tests
if [ "$RUN_INTEGRATION" = true ]; then
    echo -e "${BLUE}----------------------------------------${NC}"
    echo -e "${BLUE}Go Integration Tests${NC}"
    echo -e "${BLUE}----------------------------------------${NC}"
    echo ""
    
    if ! run_test "CLI Integration Tests" "go test ./cmd/... -tags=integration -v -timeout=60s"; then
        ALL_PASSED=false
    fi
fi

# Run bash integration tests
if [ "$RUN_BASH" = true ]; then
    echo -e "${BLUE}----------------------------------------${NC}"
    echo -e "${BLUE}Bash Integration Tests${NC}"
    echo -e "${BLUE}----------------------------------------${NC}"
    echo ""
    
    if ! run_test "Bash Integration Tests" "./test_integration.sh"; then
        ALL_PASSED=false
    fi
fi

# Summary
echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Test Summary${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

if [ "$ALL_PASSED" = true ]; then
    echo -e "${GREEN}✓ All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}✗ Some tests failed${NC}"
    exit 1
fi
