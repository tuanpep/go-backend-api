# Deployment Files

This directory contains deployment configuration files for different deployment methods.

## Files

- `go-backend-api.service` - Systemd service file (for running binary directly)
- `nginx.conf` - Nginx reverse proxy configuration (for SSL/HTTPS)

## Systemd Service Setup

If you want to run the Go binary directly (without Docker):

1. Build the binary:
   ```bash
   go build -o /usr/local/bin/go-backend-api cmd/main.go
   ```

2. Copy the service file:
   ```bash
   sudo cp deploy/go-backend-api.service /etc/systemd/system/
   ```

3. Create environment file:
   ```bash
   sudo cp env.production.example /opt/go-backend-api/.env.production
   sudo nano /opt/go-backend-api/.env.production
   ```

4. Enable and start:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable go-backend-api
   sudo systemctl start go-backend-api
   sudo systemctl status go-backend-api
   ```

## Nginx Setup

For production with SSL/TLS:

1. Install Nginx:
   ```bash
   sudo apt update
   sudo apt install nginx
   ```

2. Copy and edit configuration:
   ```bash
   sudo cp deploy/nginx.conf /etc/nginx/sites-available/go-backend-api
   sudo nano /etc/nginx/sites-available/go-backend-api
   # Update "your-domain.com" with your actual domain
   ```

3. Enable site:
   ```bash
   sudo ln -s /etc/nginx/sites-available/go-backend-api /etc/nginx/sites-enabled/
   sudo nginx -t
   sudo systemctl reload nginx
   ```

4. Install SSL certificate (Let's Encrypt):
   ```bash
   sudo apt install certbot python3-certbot-nginx
   sudo certbot --nginx -d your-domain.com -d www.your-domain.com
   ```

5. Update nginx.conf with actual certificate paths after certbot runs.

## Recommended: Docker Compose

For most use cases, using Docker Compose (see main README) is recommended as it's simpler and more portable.

