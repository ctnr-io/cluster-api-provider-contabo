#!/bin/bash
set -e

echo "Setting up Contabo credentials for e2e tests..."

# Check if required environment variables are set
required_vars=("CONTABO_CLIENT_ID" "CONTABO_CLIENT_SECRET" "CONTABO_API_USER" "CONTABO_API_PASSWORD")
missing_vars=()

for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        missing_vars+=("$var")
    fi
done

if [ ${#missing_vars[@]} -ne 0 ]; then
    echo "Error: The following required environment variables are not set:"
    for var in "${missing_vars[@]}"; do
        echo "  - $var"
    done
    echo ""
    echo "Please set these environment variables with your Contabo OAuth2 credentials:"
    echo "  export CONTABO_CLIENT_ID=\"your-client-id\""
    echo "  export CONTABO_CLIENT_SECRET=\"your-client-secret\""
    echo "  export CONTABO_API_USER=\"your-api-user\""
    echo "  export CONTABO_API_PASSWORD=\"your-api-password\""
    echo ""
    echo "You can get these credentials from the Contabo customer portal."
    exit 1
fi

echo "✓ All required Contabo credentials are set"
echo "✓ Ready to run e2e tests"

# Export variables for kustomize/envsubst
export CONTABO_CLIENT_ID
export CONTABO_CLIENT_SECRET  
export CONTABO_API_USER
export CONTABO_API_PASSWORD

echo ""
echo "Environment variables exported. You can now run:"
echo "  make test-e2e"