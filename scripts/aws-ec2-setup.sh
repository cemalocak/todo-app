#!/bin/bash

# 🚀 AWS EC2 t2.micro Setup Script for Todo App
# This script sets up a fresh Ubuntu 22.04 LTS t2.micro instance

set -e

echo "🚀 Starting AWS EC2 setup for Todo App..."

# Update system
echo "📦 Updating system packages..."
sudo yum update -y

# Install essential packages
echo "🔧 Installing essential packages..."
sudo yum install -y \
    curl \
    wget \
    git \
    unzip \
    htop \
    nano

# Install Docker
echo "🐳 Installing Docker..."
sudo yum install -y docker
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker $USER

# Install Docker Compose
echo "🐳 Installing Docker Compose..."
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Setup firewall
echo "🔥 Configuring firewall..."
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 80
sudo ufw allow 443
sudo ufw allow 8080
sudo ufw --force enable

# Configure fail2ban
echo "🛡️ Configuring fail2ban..."
sudo systemctl enable fail2ban
sudo systemctl start fail2ban

# Create app directory
echo "📁 Creating application directory..."
mkdir -p ~/todo-app
cd ~/todo-app

# Setup logging
echo "📝 Setting up logging..."
sudo mkdir -p /var/log/todo-app
sudo chown $USER:$USER /var/log/todo-app

# Create environment file
echo "⚙️ Creating environment configuration..."
cat > .env << EOF
# Production Environment
ENV=production
NODE_ENV=production

# Database
DB_PATH=/data/todos.db

# Server
PORT=8080
HOST=0.0.0.0

# Docker
BACKEND_IMAGE=ghcr.io/username/todo-app-backend:latest
FRONTEND_IMAGE=ghcr.io/username/todo-app-frontend:latest

# Monitoring
COMPOSE_PROJECT_NAME=todoapp
EOF

# Create deployment script
echo "🚀 Creating deployment script..."
cat > deploy.sh << 'EOF'
#!/bin/bash
set -e

echo "🚀 Starting deployment..."

# Load environment variables
source .env

# Login to GitHub Container Registry if token is provided
if [ ! -z "$GITHUB_TOKEN" ] && [ ! -z "$GITHUB_USER" ]; then
    echo "🔑 Logging into GitHub Container Registry..."
    echo $GITHUB_TOKEN | docker login ghcr.io -u $GITHUB_USER --password-stdin
fi

# Stop existing containers
echo "🛑 Stopping existing containers..."
docker-compose -f docker-compose.prod.yml down || true

# Pull latest images
echo "📥 Pulling latest images..."
docker pull $BACKEND_IMAGE || true
docker pull $FRONTEND_IMAGE || true

# Start new containers
echo "🚀 Starting new containers..."
docker-compose -f docker-compose.prod.yml up -d

# Wait for services to be ready
echo "⏳ Waiting for services to start..."
sleep 30

# Health check
echo "🏥 Running health checks..."
for i in {1..10}; do
    if curl -f http://localhost/api/todos > /dev/null 2>&1; then
        echo "✅ Application is healthy!"
        break
    else
        echo "⏳ Waiting for application... (attempt $i/10)"
        sleep 10
    fi
    
    if [ $i -eq 10 ]; then
        echo "❌ Health check failed!"
        docker-compose -f docker-compose.prod.yml logs
        exit 1
    fi
done

echo "🎉 Deployment completed successfully!"
echo "🌐 Application is available at: http://$(curl -s ifconfig.me)"
EOF

chmod +x deploy.sh

# Create monitoring script
echo "📊 Creating monitoring script..."
cat > monitor.sh << 'EOF'
#!/bin/bash

echo "📊 Todo App Status Report"
echo "=========================="

# System info
echo "🖥️  System Info:"
echo "   CPU: $(nproc) cores"
echo "   Memory: $(free -h | awk '/^Mem:/ { print $3"/"$2 }')"
echo "   Disk: $(df -h / | awk 'NR==2 { print $3"/"$2" ("$5" used)" }')"
echo "   Uptime: $(uptime -p)"
echo ""

# Docker status
echo "🐳 Docker Status:"
docker-compose -f docker-compose.prod.yml ps
echo ""

# Container logs (last 10 lines)
echo "📝 Recent Logs:"
echo "--- Backend ---"
docker-compose -f docker-compose.prod.yml logs --tail=5 backend
echo "--- Frontend ---"
docker-compose -f docker-compose.prod.yml logs --tail=5 frontend
echo ""

# Health check
echo "🏥 Health Check:"
if curl -f http://localhost/api/todos > /dev/null 2>&1; then
    echo "✅ Application is healthy"
else
    echo "❌ Application is not responding"
fi

# Resource usage
echo ""
echo "📈 Resource Usage:"
docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}"
EOF

chmod +x monitor.sh

# Create backup script
echo "💾 Creating backup script..."
cat > backup.sh << 'EOF'
#!/bin/bash

BACKUP_DIR="/home/ubuntu/backups"
DATE=$(date +"%Y%m%d_%H%M%S")

mkdir -p $BACKUP_DIR

echo "💾 Creating backup..."

# Backup database
docker-compose -f docker-compose.prod.yml exec -T backend cat /data/todos.db > $BACKUP_DIR/todos_${DATE}.db

# Backup configurations
tar -czf $BACKUP_DIR/config_${DATE}.tar.gz .env docker-compose.prod.yml

# Keep only last 7 backups
find $BACKUP_DIR -name "*.db" -mtime +7 -delete
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete

echo "✅ Backup completed: $BACKUP_DIR"
ls -la $BACKUP_DIR/
EOF

chmod +x backup.sh

# Setup cron job for backups
echo "⏰ Setting up automated backups..."
(crontab -l 2>/dev/null; echo "0 2 * * * /home/ubuntu/todo-app/backup.sh >> /var/log/todo-app/backup.log 2>&1") | crontab -

# Enable Docker service
echo "🔄 Enabling Docker service..."
sudo systemctl enable docker
sudo systemctl start docker

# Test Docker installation
echo "🧪 Testing Docker installation..."
docker --version
docker-compose --version

# Create systemd service for auto-restart
echo "🔄 Creating systemd service..."
sudo tee /etc/systemd/system/todo-app.service > /dev/null << EOF
[Unit]
Description=Todo App Docker Compose
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/home/ubuntu/todo-app
ExecStart=/usr/local/bin/docker-compose -f docker-compose.prod.yml up -d
ExecStop=/usr/local/bin/docker-compose -f docker-compose.prod.yml down
TimeoutStartSec=0
User=ubuntu
Group=docker

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable todo-app

echo ""
echo "🎉 AWS EC2 setup completed successfully!"
echo ""
echo "📋 Next Steps:"
echo "1. 🔑 Add your SSH key to ~/.ssh/authorized_keys"
echo "2. 🔧 Update GitHub repository URL in docker-compose.prod.yml"
echo "3. 🚀 Run: ./deploy.sh to deploy the application"
echo "4. 📊 Run: ./monitor.sh to check status"
echo "5. 💾 Run: ./backup.sh to backup data"
echo ""
echo "🌐 Your server will be available at: http://$(curl -s ifconfig.me)"
echo ""
echo "🛠️  Useful commands:"
echo "   sudo systemctl status todo-app    # Check service status"
echo "   docker-compose logs -f            # View live logs"
echo "   ./monitor.sh                      # System monitoring"
echo "   ./backup.sh                       # Manual backup"
echo ""

# Reboot message
echo "⚠️  Please reboot the system to complete Docker setup:"
echo "   sudo reboot" 