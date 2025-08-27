#!/bin/bash

# Initial Let's Encrypt certificate setup script
# Usage: ./init-letsencrypt.sh yourdomain.com your@email.com

if [ $# -ne 2 ]; then
    echo "Usage: $0 <domain> <email>"
    echo "Example: $0 circles.yourdomain.com your@email.com"
    exit 1
fi

DOMAIN=$1
EMAIL=$2

echo "### Setting up Let's Encrypt certificate for $DOMAIN..."

# Create directories
mkdir -p ./certbot/conf
mkdir -p ./certbot/www

# Replace DOMAIN placeholder in nginx config
sed -i "s/DOMAIN/$DOMAIN/g" ./nginx/nginx.conf

echo "### Starting nginx with temporary certificate..."

# Create temporary self-signed certificate for initial nginx startup
docker run --rm -v "$PWD/certbot/conf:/etc/letsencrypt" \
    certbot/certbot \
    sh -c "openssl req -x509 -nodes -newkey rsa:2048 -days 1 \
        -keyout '/etc/letsencrypt/live/$DOMAIN/privkey.pem' \
        -out '/etc/letsencrypt/live/$DOMAIN/fullchain.pem' \
        -subj '/CN=localhost' && \
        mkdir -p /etc/letsencrypt/live/$DOMAIN"

# Start services
docker compose up -d nginx

echo "### Requesting Let's Encrypt certificate for $DOMAIN..."

# Remove temporary certificate and get real one
docker run --rm -v "$PWD/certbot/conf:/etc/letsencrypt" \
    -v "$PWD/certbot/www:/var/www/certbot" \
    certbot/certbot \
    certonly --webroot -w /var/www/certbot \
    --email $EMAIL \
    --agree-tos \
    --no-eff-email \
    -d $DOMAIN

echo "### Restarting nginx with real certificate..."
docker compose restart nginx

echo "### Setting up certificate renewal..."
echo "0 12 * * * cd $PWD && docker compose run --rm certbot renew --webroot -w /var/www/certbot && docker compose restart nginx" | crontab -

echo "### Setup complete!"
echo "Your site should now be available at https://$DOMAIN"
echo "Certificate renewal is scheduled daily at noon."