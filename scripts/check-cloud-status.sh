#!/bin/bash
# Check Fly.io deployment status and health
# This script provides a comprehensive status check of the cloud deployment

set -e

# Configuration
APP_NAME="${1:-beady-backup}"

echo "=== Beady Cloud Deployment Status ==="
echo ""

# Check if flyctl is installed
if ! command -v flyctl &> /dev/null; then
    echo "ERROR: flyctl not found"
    exit 1
fi

# App information
echo "📱 App Information:"
flyctl info -a "$APP_NAME" 2>/dev/null || {
    echo "  ❌ App not found: $APP_NAME"
    exit 1
}
echo ""

# Machine status
echo "🖥️  Machine Status:"
flyctl status -a "$APP_NAME"
echo ""

# Volume information
echo "💾 Volume Information:"
flyctl volumes list -a "$APP_NAME"
echo ""

# Recent logs
echo "📋 Recent Logs (last 20 lines):"
flyctl logs -a "$APP_NAME" --lines 20
echo ""

# Database status
echo "🗄️  Database Status:"
echo "Checking database files on cloud..."
flyctl ssh console -a "$APP_NAME" -C "ls -lh /data/.beads/" 2>/dev/null || {
    echo "  ❌ Cannot access database directory"
}
echo ""

# Issue count
echo "📊 Issue Count:"
ISSUE_COUNT=$(flyctl ssh console -a "$APP_NAME" -C "cd /data/.beads && bd list --all 2>/dev/null | wc -l" 2>/dev/null || echo "0")
echo "  Total issues: $ISSUE_COUNT"
echo ""

# Health check
echo "🏥 Health Check:"
APP_URL=$(flyctl info -a "$APP_NAME" --json 2>/dev/null | jq -r '.Hostname' 2>/dev/null || echo "unknown")
if [ "$APP_URL" != "unknown" ]; then
    echo "  Testing: https://$APP_URL"
    if curl -f -s -o /dev/null -w "%{http_code}" "https://$APP_URL" 2>/dev/null | grep -q "200"; then
        echo "  ✅ App is healthy and responding"
    else
        echo "  ⚠️  App may be stopped (will auto-start on first request)"
    fi
else
    echo "  ⚠️  Cannot determine app URL"
fi
echo ""

# Quick actions
echo "🔧 Quick Actions:"
echo "  View full logs:     flyctl logs -a $APP_NAME"
echo "  SSH into machine:   flyctl ssh console -a $APP_NAME"
echo "  Open in browser:    flyctl open -a $APP_NAME"
echo "  Restart app:        flyctl apps restart $APP_NAME"
echo "  Sync to local:      bash scripts/sync-cloud-to-local.sh $APP_NAME"
echo "  Sync from local:    bash scripts/sync-local-to-cloud.sh $APP_NAME"
echo ""