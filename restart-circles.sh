#!/bin/bash

set -e

echo "🔨 Building circles-diy service..."
docker compose build circles-diy

echo "🔄 Restarting circles-diy service..."
docker compose restart circles-diy

echo "✅ circles-diy service restarted successfully!"
echo "🌐 Access the app at http://localhost or https://yourdomain.com"
echo "🎨 Concept demo available at /concept-demo/"