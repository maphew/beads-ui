#!/bin/bash
# Deploy beady to Fly.io as cloud backup endpoint
# This script automates the Fly.io deployment process

set -e

echo "=== Beady Fly.io Deployment Script ==="
echo ""

# Configuration
APP_NAME="${1:-beady-backup}"
REGION="${2:-iad}"
VOLUME_SIZE="1"

echo "Configuration:"
echo "  App Name: $APP_NAME"
echo "  Region: $REGION (Ashburn, Virginia)"
echo "  Volume Size: ${VOLUME_SIZE}GB"
echo ""

# Check if flyctl is installed
if ! command -v flyctl &> /dev/null; then
    echo "ERROR: flyctl not found"
    echo "Install it from: https://fly.io/docs/flyctl/install/"
    exit 1
fi

# Check if logged in
if ! flyctl auth whoami &> /dev/null; then
    echo "Not logged in to Fly.io"
    echo "Running: flyctl auth login"
    flyctl auth login
fi

echo "Step 1: Creating Fly.io app..."
if flyctl apps list | grep -q "$APP_NAME"; then
    echo "App $APP_NAME already exists"
else
    flyctl apps create "$APP_NAME" || {
        echo "Failed to create app. Trying with auto-generated name..."
        APP_NAME=$(flyctl apps create --json | jq -r '.Name')
        echo "Created app with name: $APP_NAME"
    }
fi

echo ""
echo "Step 2: Creating persistent volume..."
if flyctl volumes list -a "$APP_NAME" 2>/dev/null | grep -q "beady_data"; then
    echo "Volume beady_data already exists"
else
    flyctl volumes create beady_data --region "$REGION" --size "$VOLUME_SIZE" -a "$APP_NAME"
fi

echo ""
echo "Step 3: Deploying application..."
cd /workspace/beady
flyctl deploy -a "$APP_NAME"

echo ""
echo "Step 4: Verifying deployment..."
flyctl status -a "$APP_NAME"

echo ""
echo "=== Deployment Complete! ==="
echo ""
echo "App URL: https://${APP_NAME}.fly.dev"
echo ""
echo "Useful commands:"
echo "  flyctl status -a $APP_NAME          # Check app status"
echo "  flyctl logs -a $APP_NAME            # View logs"
echo "  flyctl open -a $APP_NAME            # Open in browser"
echo "  flyctl ssh console -a $APP_NAME     # SSH into machine"
echo ""
echo "Next steps:"
echo "1. Initialize database on cloud (if not syncing from local)"
echo "2. Set up Git sync for automated backup"
echo "3. Test access from mobile devices"
echo ""