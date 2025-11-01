#!/bin/bash
# Sync local Beads database to Fly.io cloud backup
# This script pushes beads.jsonl from local to cloud

set -e

# Configuration
APP_NAME="${1:-beady-backup}"
LOCAL_BEADS_DIR="${2:-$HOME/.beads}"

echo "=== Sync Local → Cloud ==="
echo "  App: $APP_NAME"
echo "  Local DB: $LOCAL_BEADS_DIR"
echo ""

# Check if flyctl is installed
if ! command -v flyctl &> /dev/null; then
    echo "ERROR: flyctl not found"
    exit 1
fi

# Check if local database exists
if [ ! -f "$LOCAL_BEADS_DIR/beads.jsonl" ]; then
    echo "ERROR: Local database not found at $LOCAL_BEADS_DIR/beads.jsonl"
    exit 1
fi

# Create backup of cloud database first
echo "Step 1: Creating backup of cloud database..."
BACKUP_FILE="/tmp/beads-cloud-backup-$(date +%Y%m%d-%H%M%S).jsonl"
flyctl ssh console -a "$APP_NAME" -C "cat /data/.beads/beads.jsonl" > "$BACKUP_FILE" 2>/dev/null || {
    echo "No existing cloud database found (this is OK for first sync)"
}

# Upload local database to cloud
echo "Step 2: Uploading local database to cloud..."
flyctl ssh console -a "$APP_NAME" -C "mkdir -p /data/.beads"
cat "$LOCAL_BEADS_DIR/beads.jsonl" | flyctl ssh console -a "$APP_NAME" -C "cat > /data/.beads/beads.jsonl"

# Copy config file if exists
if [ -f "$LOCAL_BEADS_DIR/config.json" ]; then
    echo "Step 3: Uploading config.json..."
    cat "$LOCAL_BEADS_DIR/config.json" | flyctl ssh console -a "$APP_NAME" -C "cat > /data/.beads/config.json"
fi

# Restart the app to reload data
echo "Step 4: Restarting cloud app..."
flyctl apps restart "$APP_NAME"

echo ""
echo "=== Sync Complete! ==="
echo "Local database synced to cloud"
if [ -f "$BACKUP_FILE" ]; then
    echo "Cloud backup saved to: $BACKUP_FILE"
fi
echo ""
echo "Verify at: https://${APP_NAME}.fly.dev"
