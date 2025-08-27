# circles.diy

A community platform for creatives, makers and builders.

## Production Deployment (Ubuntu VM)

### Prerequisites
- Ubuntu 24.04 LTS VM with public IP
- Docker and Docker Compose installed
- Domain pointing to your VM's IP

### Quick Setup

1. **Clone and prepare:**
   ```bash
   git clone <your-repo>
   cd circles.diy
   ```

2. **Initialize SSL certificate:**
   ```bash
   ./init-letsencrypt.sh yourdomain.com your@email.com
   ```
   
   This script will:
   - Update nginx config with your domain
   - Start nginx with a temporary certificate
   - Request a real Let's Encrypt certificate
   - Restart nginx with the real certificate
   - **Start your circles.diy app**

Your site will be available at `https://yourdomain.com`

### Manual Startup (if needed)
```bash
docker compose up -d
```

### Manual SSL Renewal
```bash
docker compose --profile renewal run --rm certbot
docker compose restart nginx
```

### Security Features
- ✅ HTTPS with Let's Encrypt
- ✅ Rate limiting (10 req/min general, 5 req/min feedback)
- ✅ Security headers (HSTS, CSP, XSS protection)
- ✅ Input validation and sanitization
- ✅ Non-root container execution
- ✅ CSRF protection

### Local Development
```bash
docker compose up circles-diy
# Access at http://localhost:6969
```

### Monitoring
- **View logs:** `docker compose logs -f`
- **Check certificates:** `docker compose exec nginx nginx -t`
- **Feedback storage:** `./data/feedback.txt`

### Firewall Setup
```bash
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp  
sudo ufw allow 443/tcp
sudo ufw enable
```