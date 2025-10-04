#!/bin/bash

# Exit on any error
set -e

# Check if environment variables are set
if [ -z "$CI_BASE_URL" ]; then
    echo "Error: CI_BASE_URL environment variable is not set"
    exit 1
fi

if [ -z "$CI_API_KEY" ]; then
    echo "Error: CI_API_KEY environment variable is not set"
    exit 1
fi

# Install newman if not already installed
if ! command -v newman &> /dev/null; then
    echo "Installing newman..."
    npm install -g newman
fi

# Run the tests
echo "Running Postman tests against $CI_BASE_URL..."
newman run ./postman/RemindMeBot.postman_collection.json \
    -e ./postman/RemindMeBot.ci.postman_environment.json \
    --env-var "CI_BASE_URL=$CI_BASE_URL" \
    --env-var "CI_API_KEY=$CI_API_KEY" \
    --reporters cli,junit \
    --reporter-junit-export test-results.xml