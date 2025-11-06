# Fly.io Deployment Guide for Beady
## Cloud Backup/Sync Endpoint Setup

This guide covers deploying beady to Fly.io as a cloud backup and sync endpoint for your local-first beady deployment.

---

## Prerequisites

### 1. Install Fly.io CLI (flyctl)
```bash
# Linux
curl -L https://fly.io/install.sh | sh

# macOS
brew install flyctl

# Verify installation
flyctl version
```

### 2. Sign up / Log in to Fly.io
```bash
# If you have an account
flyctl auth login

# If you need to sign up
flyctl auth signup
```

### 3. Verify Free Tier Eligibility
```bash
# Check your account status
flyctl status

# This deployment is configured to stay within the free tier:
# - 1 shared-cpu-1x VM (256MB RAM)
# - 1GB persistent volume
# - Auto-stop when idle (no ongoing costs)
```

---

## Deployment Steps

### Step 1: Navigate to beady directory
```bash
cd /workspace/beady
```

### Step 2: Create Fly.io App
```bash
# Create app (change app name if needed)
flyctl apps create beady-backup

# Or let Fly.io generate a unique name
flyctl apps create
```

### Step 3: Create Persistent Volume
```bash
# Create 1GB volume in your primary region
flyctl volumes create beady_data --region iad --size 1

# Verify volume creation
flyctl volumes list
```

### Step 4: Deploy Application
```bash
# Deploy to Fly.io
flyctl deploy

# This will:
# 1. Build Docker image from Dockerfile
# 2. Push image to Fly.io registry
# 3. Create and start the Machine
# 4. Mount the volume to /data
```

### Step 5: Verify Deployment
```bash
# Check app status
flyctl status

# View app URL
flyctl info

# Access logs
flyctl logs

# Open in browser
flyctl open
```

---

## Configuration Details

### App Settings (fly.toml)
- **App Name**: beady-backup (customize as needed)
- **Region**: iad (Ashburn, Virginia)
- **VM Size**: shared-cpu-1x (256MB RAM, 1 vCPU)
- **Port**: 8080 (internal), 443 (external HTTPS)
- **Auto-scaling**: Stops when idle, starts on request

### Volume Configuration
- **Name**: beady_data
- **Size**: 1GB (free tier)
- **Mount Point**: /data
- **Database Path**: /data/.beads/
- **Snapshots**: Daily, retained for 5 days
- **Auto-extend**: Enabled at 80% usage

### Cost Optimization
This configuration is designed to minimize costs:
- ✅ Auto-stop when idle (no traffic)
- ✅ Auto-start on new requests
- ✅ No always-on machines (min_machines_running = 0)
- ✅ 1GB volume (within free tier)
- ✅ Shared CPU (cost-effective)

**Expected Monthly Cost**: $0-5 with light usage

---

## Database Initialization

### Option 1: Initialize Fresh Database
After first deployment, the cloud instance will have an empty database. Initialize it:

```bash
# SSH into the Fly.io machine
flyctl ssh console

# Inside the machine
cd /data/.beads
bd init --prefix retirement-project

# Exit SSH
exit
```

### Option 2: Sync from Local Database
See the "Sync Strategies" section below to push your local database to the cloud.

---

## Sync Strategies

### 1. Git-Based Sync (Recommended)

The Beads database uses `beads.jsonl` as the source of truth, which is perfect for Git sync.

#### Setup Git Sync
```bash
# On local machine
cd ~/.beads
git init
git add beads.jsonl
git commit -m "Initial beady database"

# Add remote (GitHub, GitLab, etc.)
git remote add origin YOUR_GIT_REPO_URL
git push -u origin main
```

#### Sync Local → Cloud
```bash
# On local machine
cd ~/.beads
git add beads.jsonl
git commit -m "Update beads database"
git push

# On Fly.io machine (via SSH)
flyctl ssh console
cd /data/.beads
git pull
exit

# Restart beady to reload data
flyctl apps restart beady-backup
```

#### Sync Cloud → Local
```bash
# On Fly.io machine (via SSH)
flyctl ssh console
cd /data/.beads
git add beads.jsonl
git commit -m "Update from cloud"
git push
exit

# On local machine
cd ~/.beads
git pull
```

### 2. Manual File Sync

#### Push Local Database to Cloud
```bash
# Copy beads.jsonl to cloud
flyctl ssh console -C "mkdir -p /data/.beads"
flyctl ssh sftp shell

# In SFTP shell
put ~/.beads/beads.jsonl /data/.beads/beads.jsonl
bye

# Restart beady
flyctl apps restart beady-backup
```

#### Pull Cloud Database to Local
```bash
# Download beads.jsonl from cloud
flyctl ssh sftp shell

# In SFTP shell
get /data/.beads/beads.jsonl ~/.beads/beads.jsonl
bye

# Restart local beady
pkill beady
cd ~/.beads && beady --port 8081 &
```

---

## Testing the Deployment

### 1. Test Web UI Access
```bash
# Get your app URL
flyctl info

# Visit the URL in browser (will auto-start the machine)
# Example: https://beady-backup.fly.dev
```

### 2. Test Auto-Stop/Start
```bash
# The machine should stop after ~5 minutes of inactivity
flyctl status

# Visit the URL again - machine should auto-start
```

### 3. Test Data Persistence
```bash
# Create an issue via the web UI
# Restart the app
flyctl apps restart beady-backup

# Verify the issue still exists (data persisted to volume)
```

---

## Backup and Recovery

### Volume Snapshots
Fly.io automatically creates daily snapshots of your volume:

```bash
# List snapshots
flyctl volumes snapshots list beady_data

# Create manual snapshot
flyctl volumes snapshots create beady_data

# Restore from snapshot
flyctl volumes create beady_data_restore --snapshot-id SNAPSHOT_ID
```

### Database Backup
```bash
# Download beads.jsonl for local backup
flyctl ssh sftp shell
get /data/.beads/beads.jsonl ~/backups/beads-backup-$(date +%Y%m%d).jsonl
bye
```

---

## Monitoring and Maintenance

### View Logs
```bash
# Real-time logs
flyctl logs

# Specific time range
flyctl logs --since 1h
```

### Resource Usage
```bash
# App metrics
flyctl status

# Volume usage
flyctl volumes list
```

### Scaling (if needed)
```bash
# Increase VM size
flyctl scale vm shared-cpu-2x

# Increase memory
flyctl scale memory 512

# Extend volume
flyctl volumes extend beady_data --size 2
```

---

## Troubleshooting

### Machine Won't Start
```bash
# Check logs
flyctl logs

# Restart machine
flyctl apps restart beady-backup

# SSH and debug
flyctl ssh console
```

### Database Not Persisting
```bash
# Verify volume mount
flyctl ssh console
ls -la /data/.beads/

# Check volume status
flyctl volumes list
```

### Out of Memory
```bash
# Increase memory allocation
flyctl scale memory 512
```

---

## Security Considerations

### 1. Make App Private (Optional)
If you only need cloud backup, not public access:

```bash
# Remove public IP (private only)
flyctl ips list
flyctl ips release <PUBLIC_IP>
```

### 2. Add Authentication (Future Enhancement)
Consider adding basic auth or OAuth if exposing publicly.

### 3. Git Repository Access
- Use private Git repositories for sensitive data
- Configure SSH keys or tokens for automated sync

---

## Next Steps

After successful deployment:

1. ✅ Test cloud access from mobile devices
2. ✅ Set up automated Git sync (cron job or GitHub Actions)
3. ✅ Configure backup rotation strategy
4. ✅ Test failover scenarios (local down, use cloud)
5. ✅ Document restore procedures

---

## Cleanup (if needed)

To remove the Fly.io deployment:

```bash
# Delete the app and all resources
flyctl apps destroy beady-backup

# This will also delete the volume and all data
```

---

## Cost Summary

**Free Tier Limits** (as of 2025):
- Up to $5/month in compute usage
- 3GB persistent volume storage (1GB used)
- Outbound data transfer limits

**This Deployment**:
- ~$0-2/month with auto-stop enabled
- Additional costs only if high traffic

**To stay free**: Ensure `auto_stop_machines = 'stop'` and `min_machines_running = 0` in fly.toml.
