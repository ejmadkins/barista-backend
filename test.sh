#!/bin/bash
# Run go tests

BACKEND_DIR="$( dirname -- "$0"; )";
echo "Running all tests found in $( pwd; )";
echo "${BACKEND_DIR}";

go test "${BACKEND_DIR}" -v;