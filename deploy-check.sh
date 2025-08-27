#!/bin/bash

# Deployment verification script
# Usage: ./deploy-check.sh yourdomain.com

DOMAIN=${1:-"localhost"}
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🔍 circles.diy Deployment Check for $DOMAIN"
echo "=============================================="

# Check if running in Docker
check_containers() {
    echo -n "📦 Docker containers: "
    if docker compose ps --services --filter "status=running" | grep -q "circles-diy\|nginx"; then
        RUNNING=$(docker compose ps --services --filter "status=running" | wc -l)
        echo -e "${GREEN}✅ $RUNNING services running${NC}"
        docker compose ps
    else
        echo -e "${RED}❌ No services running${NC}"
        return 1
    fi
}

# Check HTTP connectivity
check_http() {
    echo -n "🌐 HTTP connectivity: "
    if curl -sf -m 10 "http://$DOMAIN/" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ HTTP accessible${NC}"
    else
        echo -e "${RED}❌ HTTP not accessible${NC}"
        return 1
    fi
}

# Check HTTPS connectivity
check_https() {
    echo -n "🔒 HTTPS connectivity: "
    if curl -sf -m 10 "https://$DOMAIN/" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ HTTPS accessible${NC}"
    else
        echo -e "${YELLOW}⚠️  HTTPS not accessible${NC}"
        return 1
    fi
}

# Check SSL certificate
check_ssl_cert() {
    echo -n "📄 SSL certificate: "
    if timeout 10 openssl s_client -connect "$DOMAIN:443" -servername "$DOMAIN" </dev/null 2>/dev/null | grep -q "Verify return code: 0"; then
        EXPIRY=$(timeout 10 openssl s_client -connect "$DOMAIN:443" -servername "$DOMAIN" </dev/null 2>/dev/null | openssl x509 -noout -dates 2>/dev/null | grep "notAfter" | cut -d= -f2)
        echo -e "${GREEN}✅ Valid (expires: $EXPIRY)${NC}"
    else
        echo -e "${YELLOW}⚠️  Certificate issues detected${NC}"
        return 1
    fi
}

# Check application functionality
check_app_functionality() {
    echo -n "⚙️  App functionality: "
    RESPONSE=$(curl -sf -m 10 "http://$DOMAIN/" 2>/dev/null || curl -sf -m 10 "https://$DOMAIN/" 2>/dev/null)
    if echo "$RESPONSE" | grep -q "circles.diy"; then
        echo -e "${GREEN}✅ App responding correctly${NC}"
    else
        echo -e "${RED}❌ App not responding correctly${NC}"
        return 1
    fi
}

# Check feedback form
check_feedback_form() {
    echo -n "📝 Feedback form: "
    RESPONSE=$(curl -sf -m 10 "http://$DOMAIN/" 2>/dev/null || curl -sf -m 10 "https://$DOMAIN/" 2>/dev/null)
    if echo "$RESPONSE" | grep -q "feedback" && echo "$RESPONSE" | grep -q "csrf_token"; then
        echo -e "${GREEN}✅ Form present with CSRF protection${NC}"
    else
        echo -e "${YELLOW}⚠️  Form issues detected${NC}"
        return 1
    fi
}

# Check security headers
check_security_headers() {
    echo -n "🛡️  Security headers: "
    HEADERS=$(curl -sI -m 10 "https://$DOMAIN/" 2>/dev/null || curl -sI -m 10 "http://$DOMAIN/" 2>/dev/null)
    
    MISSING=()
    echo "$HEADERS" | grep -qi "X-Frame-Options" || MISSING+=("X-Frame-Options")
    echo "$HEADERS" | grep -qi "X-Content-Type-Options" || MISSING+=("X-Content-Type-Options")
    echo "$HEADERS" | grep -qi "Content-Security-Policy" || MISSING+=("CSP")
    
    if [ ${#MISSING[@]} -eq 0 ]; then
        echo -e "${GREEN}✅ All security headers present${NC}"
    else
        echo -e "${YELLOW}⚠️  Missing: ${MISSING[*]}${NC}"
        return 1
    fi
}

# Check file permissions
check_file_permissions() {
    echo -n "🔐 File permissions: "
    if [ -f "./data/feedback.txt" ]; then
        PERMS=$(stat -c "%a" "./data/feedback.txt" 2>/dev/null || echo "unknown")
        if [ "$PERMS" = "644" ] || [ "$PERMS" = "640" ]; then
            echo -e "${GREEN}✅ Proper permissions ($PERMS)${NC}"
        else
            echo -e "${YELLOW}⚠️  Unusual permissions ($PERMS)${NC}"
        fi
    else
        echo -e "${YELLOW}⚠️  No feedback file yet${NC}"
    fi
}

# Run all checks
echo ""
FAILED=0

check_containers || FAILED=$((FAILED+1))
check_http || FAILED=$((FAILED+1))
check_https || FAILED=$((FAILED+1))
check_ssl_cert || FAILED=$((FAILED+1))
check_app_functionality || FAILED=$((FAILED+1))
check_feedback_form || FAILED=$((FAILED+1))
check_security_headers || FAILED=$((FAILED+1))
check_file_permissions

echo ""
echo "=============================================="
if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 All checks passed! Deployment looks good.${NC}"
    exit 0
elif [ $FAILED -le 2 ]; then
    echo -e "${YELLOW}⚠️  $FAILED minor issues detected. Deployment mostly working.${NC}"
    exit 0
else
    echo -e "${RED}❌ $FAILED issues detected. Deployment needs attention.${NC}"
    exit 1
fi