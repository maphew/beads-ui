#!/bin/bash
# Sync Fly.io cloud database to local Beads backup
# This script pulls beads.jsonl from cloud to local

set -e

# Configuration
APP_NAME="${1:-beady-backup}"
LOCAL_BEADS_DIR="${2:-$HOME/.beads}"

echo "=== Sync Cloud → Local ==="
echo "  App: $APP_NAME"
echo "  Local DB: $LOCAL_BEADS_DIR"
echo ""

# Check if flyctl is installed
if ! command -v flyctl &> /dev/null; then
    echo "ERROR: flyctl not found"
    exit 1
fi

# Create backup of local database first
if [ -f "$LOCAL_BEADS_DIR/beads.jsonl" ]; then
    echo "Step 1: Creating backup of local database..."
    BACKUP_FILE="$LOCAL_BEADS_DIR/beads-local-backup-$(date +%Y%m%d-%H%M%S).jsonl"
    cp "$LOCAL_BEADS_DIR/beads.jsonl" "$BACKUP_FILE"
    echo "Local backup saved to: $BACKUP_FILE"
fi

# Create local directory if not exists
mkdir -p "$LOCAL_BEADS_DIR"

# Download cloud database
echo "Step 2: Downloading cloud database..."
flyctl ssh console -a "$APP_NAME" -C "cat /data/.beads/beads.jsonl" > "$LOCAL_BEADS_DIR/beads.jsonl"

# Download config file
echo "Step 3: Downloading config.json..."
flyctl ssh console -a "$APP_NAME" -C "cat /data/.beads/config.json" > "$LOCAL_BEADS_DIR/config.json" 2>/dev/null || {
    echo "No config.json found on cloud (this is OK)"
}

# Restart local beady if running
echo "Step 4: Restarting local beady (if running)..."
pkill beady 2>/dev/null || echo "Local beady not running"

echo ""
echo "=== Sync Complete! ==="
echo "Cloud database synced to local"
echo "Local database: $LOCAL_BEADS_DIR/beads.jsonl"
if [ -f "$BACKUP_FILE" ]; then
    echo "Previous local backup: $BACKUP_FILE"
fi
echo ""
echo "To start local beady:"
echo "  cd $LOCAL_BEADS_DIR && beady --port 8081"
