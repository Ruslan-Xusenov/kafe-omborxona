#!/bin/bash
set -e

echo "========================================="
echo "  Kafe Omborxona Production Deployment"
echo "========================================="

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

PROJECT_DIR="/opt/kafe-omborxona"
REPO_URL="https://github.com/Ruslan-Xusenov/kafe-omborxona.git"
SERVER_IP="157.180.118.115"

print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

echo ""
echo "1. Updating system packages..."
apt-get update
apt-get upgrade -y
print_status "System updated"

echo ""
echo "2. Installing required packages..."
apt-get install -y git curl nginx certbot python3-certbot-nginx
print_status "Packages installed"

echo ""
echo "3. Installing Docker & Docker Compose..."
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com -o get-docker.sh
    sh get-docker.sh
    rm get-docker.sh
fi

if ! command -v docker-compose &> /dev/null; then
    curl -L "https://github.com/docker/compose/releases/download/v2.24.5/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
fi
print_status "Docker & Docker Compose ready"

echo ""
echo "4. Cloning repository..."
if [ -d "$PROJECT_DIR" ]; then
    echo "Directory exists. Pulling latest changes..."
    cd $PROJECT_DIR
    git reset --hard
    git pull
else
    git clone $REPO_URL $PROJECT_DIR
    cd $PROJECT_DIR
fi
print_status "Repository cloned"

echo ""
echo "5. Setting up Environment Variables (.env)..."
mkdir -p $PROJECT_DIR/backend
mkdir -p $PROJECT_DIR/frontend

cat <<EOF > $PROJECT_DIR/backend/.env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=kafe_omborxona
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=kafe_omborxona
JWT_SECRET=kafe-omborxona-super-secret-key-production-2026
SERVER_PORT=8080
FRONTEND_URL=http://$SERVER_IP
TELEGRAM_BOT_TOKEN=8988171225:AAGJDA172tE9bzMyRDf0Uxg93yN2ndT8vBI
EOF

cat <<EOF > $PROJECT_DIR/frontend/.env.local
NEXT_PUBLIC_API_URL=/api
EOF
print_status "Environment variables generated"

echo ""
echo "6. Starting Docker Containers..."
cd $PROJECT_DIR
docker-compose down
docker-compose build db backend frontend
docker-compose up -d db backend frontend
print_status "Containers are running"

echo ""
echo "7. Configuring Host Nginx..."
cat <<EOF > /etc/nginx/sites-available/kafe-omborxona
server {
    listen 80;
    server_name $SERVER_IP;

    # Backend API xizmati
    location /api/ {
        proxy_pass http://localhost:8080/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_cache_bypass \$http_upgrade;
    }

    # Frontend Next.js xizmati
    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_cache_bypass \$http_upgrade;
    }
}
EOF

ln -sf /etc/nginx/sites-available/kafe-omborxona /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default
nginx -t
systemctl restart nginx
systemctl enable nginx
print_status "Nginx configured"

echo ""
echo "8. Configuring Firewall..."
if command -v ufw &> /dev/null; then
    ufw allow 22/tcp
    ufw allow 80/tcp
    ufw allow 443/tcp
    ufw --force enable
    print_status "Firewall configured"
fi

echo ""
echo "========================================="
echo -e "${GREEN}  Deployment Complete!${NC}"
echo "========================================="
echo "Saytingiz quyidagi manzilda ishga tushdi:"
echo "http://$SERVER_IP"
echo ""
echo "ESLATMA SSL (HTTPS) uchun domen ulashingiz kerak."
echo "Domen ulagandan so'ng serverda shu komandani yozasiz:"
echo "sudo certbot --nginx -d sizning-domeningiz.uz"
echo ""
print_status "Barcha ishlar muvaffaqiyatli yakunlandi!"