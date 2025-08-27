# circles.diy

A community platform for creatives, makers and builders.

Read the manifesto: https://circles.diy/

## Production Deployment (Ubuntu VM)

### Prerequisites
- Ubuntu 24.04 LTS VM with public IP
- Docker and Docker Compose installed
- Domain pointing to your VM's IP
- Ports 80 and 443 open in firewall

### Step-by-Step Setup

1. **Install Docker (if needed):**
   ```bash
   sudo apt update
   sudo apt install docker.io docker-compose-v2 -y
   sudo usermod -aG docker $USER
   # Log out and back in
   ```

2. **Configure firewall:**
   ```bash
   sudo ufw allow 22/tcp   # SSH
   sudo ufw allow 80/tcp   # HTTP
   sudo ufw allow 443/tcp  # HTTPS
   sudo ufw enable
   ```

3. **Clone and deploy:**
   ```bash
   git clone <your-repo>
   cd circles.diy
   ./init-letsencrypt.sh yourdomain.com your@email.com
   ```

4. **Verify deployment:**
   ```bash
   ./deploy-check.sh yourdomain.com
   ```

### What the Setup Does
The `init-letsencrypt.sh` script:
- ✅ Validates your domain and email
- ✅ Checks domain DNS resolution  
- ✅ Starts services with HTTP-only nginx
- ✅ Requests Let's Encrypt SSL certificate
- ✅ Switches to HTTPS configuration
- ✅ Tests the final deployment
- ✅ Sets up automatic certificate renewal

### Troubleshooting

**Certificate request fails:**
```bash
# Check domain resolution
dig yourdomain.com A

# Check HTTP accessibility
curl -I http://yourdomain.com

# View nginx logs
docker compose logs nginx
```

**Service not starting:**
```bash
# Check all container status
docker compose ps

# View app logs
docker compose logs circles-diy

# Restart services
docker compose restart
```

### Manual Operations

**Manual SSL renewal:**
```bash
docker compose --profile renewal run --rm certbot
docker compose restart nginx
```

**View logs:**
```bash
docker compose logs -f
```

**Update and redeploy:**
```bash
git pull
docker compose build --no-cache
docker compose up -d
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