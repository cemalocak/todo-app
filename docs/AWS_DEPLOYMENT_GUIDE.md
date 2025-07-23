# 🚀 AWS EC2 Deployment Guide

Complete guide to deploy Todo App on AWS EC2 t2.micro (Free Tier) with CI/CD Pipeline.

## 📋 Prerequisites

- ✅ AWS Account (Free Tier eligible)
- ✅ GitHub Account
- ✅ Domain name (optional)
- ✅ SSH Key pair

---

## 🔧 Step 1: AWS EC2 Setup

### 1.1 Launch EC2 Instance

1. **Login to AWS Console** → EC2 Dashboard
2. **Launch Instance**:
   - **Name**: `todo-app-production`
   - **AMI**: Ubuntu 22.04 LTS (Free tier eligible)
   - **Instance Type**: `t2.micro` (Free tier eligible)
   - **Key Pair**: Create or select existing SSH key
   - **Security Group**:
     ```
     SSH (22)     - Your IP only
     HTTP (80)    - 0.0.0.0/0
     HTTPS (443)  - 0.0.0.0/0
     Custom (8080) - 0.0.0.0/0 (for API)
     ```
   - **Storage**: 30 GB gp2 (Free tier eligible)

3. **Launch Instance** and wait for it to be running

### 1.2 Connect to Instance

```bash
# SSH to your instance
ssh -i "your-key.pem" ubuntu@YOUR_EC2_PUBLIC_IP

# Make setup script executable and run
wget https://raw.githubusercontent.com/YOUR_USERNAME/YOUR_REPO/main/scripts/aws-ec2-setup.sh
chmod +x aws-ec2-setup.sh
./aws-ec2-setup.sh

# Reboot to complete Docker setup
sudo reboot
```

---

## 🔐 Step 2: GitHub Repository Setup

### 2.1 Create Repository

1. Create new GitHub repository: `todo-app`
2. Push your code:

```bash
git init
git add .
git commit -m "Initial commit: Complete Todo App with CI/CD"
git branch -M main
git remote add origin https://github.com/YOUR_USERNAME/todo-app.git
git push -u origin main

# Create develop branch for testing
git checkout -b develop
git push -u origin develop
```

### 2.2 Configure GitHub Secrets

Go to **Repository Settings → Secrets and Variables → Actions**:

#### Required Secrets:

```yaml
# AWS Configuration
AWS_ACCESS_KEY_ID: "AKIA..."          # Your AWS Access Key
AWS_SECRET_ACCESS_KEY: "..."          # Your AWS Secret Key  
AWS_REGION: "us-east-1"               # Your preferred region

# EC2 Configuration
EC2_HOST: "3.84.123.45"               # Your EC2 public IP
EC2_SSH_KEY: |                        # Your private SSH key (full content)
  -----BEGIN RSA PRIVATE KEY-----
  MIIEpAIBAAKCAQEA...
  -----END RSA PRIVATE KEY-----

# Environment URLs
TEST_API_URL: "http://YOUR_EC2_IP:8080"     # Test environment
PROD_API_URL: "http://YOUR_EC2_IP"          # Production environment

# GitHub Container Registry (Automatic)
GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}   # Auto-provided by GitHub
```

#### Optional Secrets:
```yaml
# Notifications (optional)
SLACK_WEBHOOK_URL: "https://hooks.slack.com/..."
DISCORD_WEBHOOK_URL: "https://discord.com/api/webhooks/..."

# Monitoring (optional)
SENTRY_DSN: "https://..."
```

---

## 🌐 Step 3: Repository Configuration

### 3.1 Update Docker Compose

Edit `docker-compose.prod.yml` and replace:
```yaml
services:
  backend:
    image: ghcr.io/YOUR_USERNAME/todo-app-backend:latest
  frontend:
    image: ghcr.io/YOUR_USERNAME/todo-app-frontend:latest
```

### 3.2 Enable GitHub Container Registry

1. Go to **Profile → Settings → Developer Settings → Personal Access Tokens**
2. Generate token with `write:packages` permission
3. Set as `GITHUB_TOKEN` secret (usually auto-provided)

---

## 🚀 Step 4: Deploy & Test

### 4.1 Initial Deployment

```bash
# SSH to your EC2 instance
ssh -i "your-key.pem" ubuntu@YOUR_EC2_IP

# Navigate to app directory
cd ~/todo-app

# Manual first deployment
./deploy.sh

# Check status
./monitor.sh
```

### 4.2 Test Application

```bash
# Test API
curl http://YOUR_EC2_IP:8080/api/todos

# Test Frontend
curl http://YOUR_EC2_IP

# Monitor logs
docker-compose -f docker-compose.prod.yml logs -f
```

---

## 🔄 Step 5: CI/CD Pipeline

### 5.1 Branch Strategy

```bash
main branch    → Production deployment
develop branch → Test environment deployment
feature/*      → PR checks only
```

### 5.2 Trigger Deployments

#### Test Deployment:
```bash
git checkout develop
git add .
git commit -m "feat: new feature"
git push origin develop
```

#### Production Deployment:
```bash
git checkout main
git merge develop
git push origin main
```

### 5.3 Pipeline Stages

```mermaid
graph LR
    A[Push Code] --> B[Run Tests]
    B --> C[Build Images]
    C --> D[Deploy Test]
    D --> E[API Tests]
    D --> F[E2E Tests]
    E --> G[Deploy Prod]
    F --> G
    G --> H[Health Check]
```

---

## 📊 Step 6: Monitoring & Maintenance

### 6.1 Monitoring Commands

```bash
# SSH to EC2
ssh -i "your-key.pem" ubuntu@YOUR_EC2_IP

# Check system status
./monitor.sh

# View logs
docker-compose -f docker-compose.prod.yml logs -f

# Check disk usage
df -h

# Check memory usage
free -h

# Check running containers
docker ps
```

### 6.2 Backup & Recovery

```bash
# Manual backup
./backup.sh

# List backups
ls -la ~/backups/

# Restore from backup (if needed)
docker-compose -f docker-compose.prod.yml down
# Copy backup files to volume location
docker-compose -f docker-compose.prod.yml up -d
```

### 6.3 SSL Certificate (Optional)

```bash
# Install Certbot
sudo apt-get install certbot python3-certbot-nginx

# Get certificate (replace YOUR_DOMAIN)
sudo certbot --nginx -d YOUR_DOMAIN

# Test renewal
sudo certbot renew --dry-run
```

---

## 🚨 Troubleshooting

### Common Issues:

#### 1. Container Won't Start
```bash
# Check logs
docker-compose -f docker-compose.prod.yml logs backend
docker-compose -f docker-compose.prod.yml logs frontend

# Restart services
docker-compose -f docker-compose.prod.yml restart
```

#### 2. Out of Disk Space
```bash
# Clean unused images
docker system prune -a

# Check disk usage
df -h
du -sh ~/todo-app/*
```

#### 3. Memory Issues (t2.micro has 1GB RAM)
```bash
# Check memory usage
free -h

# Restart containers to free memory
docker-compose -f docker-compose.prod.yml restart
```

#### 4. SSH Connection Issues
```bash
# Check security group allows SSH from your IP
# Verify SSH key permissions (chmod 600 your-key.pem)
# Try connecting with verbose output
ssh -v -i "your-key.pem" ubuntu@YOUR_EC2_IP
```

#### 5. GitHub Actions Failing
```bash
# Check secrets are set correctly
# Verify SSH key format (include headers/footers)
# Check EC2 instance is running
# Verify security groups allow connections
```

---

## 💰 Cost Optimization

### Free Tier Limits:
- ✅ **EC2 t2.micro**: 750 hours/month (1 instance always-on)
- ✅ **EBS Storage**: 30 GB/month
- ✅ **Data Transfer**: 1 GB/month outbound

### Best Practices:
- ✅ Monitor usage in AWS Billing Dashboard
- ✅ Set up billing alerts
- ✅ Stop instance when not needed (development)
- ✅ Use CloudWatch free tier for monitoring

---

## 🔗 Useful Links

- [AWS Free Tier](https://aws.amazon.com/free/)
- [GitHub Actions Docs](https://docs.github.com/en/actions)
- [Docker Compose Docs](https://docs.docker.com/compose/)
- [Ubuntu Server Guide](https://ubuntu.com/server/docs)

---

## 📞 Support

If you encounter issues:

1. 📖 Check this guide first
2. 🔍 Search GitHub Issues
3. 📝 Create new issue with:
   - Error logs
   - Steps to reproduce
   - Environment details

---

**🎉 Your Todo App is now running on AWS with full CI/CD pipeline!**

**URLs:**
- 🌐 **Application**: `http://YOUR_EC2_IP`
- 🔧 **API**: `http://YOUR_EC2_IP:8080/api/todos`
- 📊 **Monitoring**: `ssh ubuntu@YOUR_EC2_IP && ./monitor.sh` 