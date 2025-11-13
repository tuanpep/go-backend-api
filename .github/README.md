# GitHub Actions CI/CD Setup

This repository includes GitHub Actions workflows for continuous integration and deployment.

## Workflows

### 1. CI Workflow (`.github/workflows/ci.yml`)

Runs on every push and pull request to `main` and `develop` branches.

**Jobs:**
- **Test**: Runs Go tests with PostgreSQL service
- **Build**: Builds the Go binary
- **Docker Build**: Builds and tests Docker image
- **Lint**: Runs golangci-lint for code quality

**Features:**
- Parallel job execution
- Code coverage reporting (Codecov)
- Docker image validation
- Automated testing

### 2. Deploy Workflow (`.github/workflows/deploy.yml`)

Deploys the application to your VM when code is pushed to `main` branch.

**Triggers:**
- Automatic: Push to `main` branch
- Manual: Workflow dispatch with environment selection

**Steps:**
1. Checks out code
2. Sets up SSH connection to VM
3. Syncs code to VM
4. Runs deployment script
5. Verifies deployment

### 3. Docker Build & Push (`.github/workflows/docker-build-push.yml`)

Builds and pushes Docker images to Docker Hub.

**Triggers:**
- Push to `main` branch
- Tag push (e.g., `v1.0.0`)
- Manual workflow dispatch

**Features:**
- Multi-platform builds (amd64, arm64)
- Automatic tagging
- Docker Hub integration

## Setup Instructions

### 1. Configure GitHub Secrets

Go to your repository → Settings → Secrets and variables → Actions → New repository secret

#### Required Secrets for Deployment:

| Secret Name | Description | Example |
|------------|-------------|---------|
| `VM_HOST` | Your VM's IP address or hostname | `192.168.1.100` or `api.example.com` |
| `VM_USER` | SSH username for VM | `deploy` or `ubuntu` |
| `VM_SSH_PRIVATE_KEY` | SSH private key for VM access | Contents of `~/.ssh/id_rsa` |
| `API_PORT` | (Optional) API port on VM | `8080` (default) |

#### Optional Secrets for Docker Hub:

| Secret Name | Description |
|------------|-------------|
| `DOCKER_USERNAME` | Docker Hub username |
| `DOCKER_PASSWORD` | Docker Hub password or access token |

### 2. Generate SSH Key Pair

On your local machine:

```bash
# Generate SSH key (if you don't have one)
ssh-keygen -t ed25519 -C "github-actions" -f ~/.ssh/github_actions_deploy

# Copy public key to VM
ssh-copy-id -i ~/.ssh/github_actions_deploy.pub user@your-vm-host

# Display private key (copy this to GitHub secret)
cat ~/.ssh/github_actions_deploy
```

**Important:** Never commit the private key to the repository!

### 3. Prepare Your VM

On your VM, ensure:

```bash
# Install Docker and Docker Compose
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Install Docker Compose plugin
sudo apt install docker-compose-plugin

# Create deployment directory
sudo mkdir -p /opt/go-backend-api
sudo chown $USER:$USER /opt/go-backend-api

# Create .env.production file
cd /opt/go-backend-api
cp env.production.example .env.production
nano .env.production  # Edit with your production values
```

### 4. Configure GitHub Environments (Optional)

For better security and environment-specific deployments:

1. Go to Settings → Environments
2. Create `production` and `staging` environments
3. Add environment-specific secrets if needed
4. Configure protection rules (required reviewers, etc.)

### 5. Test the Workflows

#### Test CI:
```bash
# Create a test branch
git checkout -b test-ci

# Make a small change and push
git add .
git commit -m "Test CI workflow"
git push origin test-ci

# Create a pull request to main
```

#### Test Deployment:
```bash
# Merge to main (or push directly)
git checkout main
git merge test-ci
git push origin main

# Or trigger manually:
# Go to Actions → Deploy → Run workflow
```

## Workflow Status Badge

Add this to your README.md to show CI status:

```markdown
![CI](https://github.com/username/go-backend-api/workflows/CI/badge.svg)
![Deploy](https://github.com/username/go-backend-api/workflows/Deploy/badge.svg)
```

## Troubleshooting

### Deployment Fails

1. **SSH Connection Issues:**
   ```bash
   # Test SSH connection manually
   ssh -i ~/.ssh/github_actions_deploy user@your-vm-host
   ```

2. **Permission Issues:**
   ```bash
   # On VM, check permissions
   ls -la /opt/go-backend-api
   chmod +x /opt/go-backend-api/scripts/deploy.sh
   ```

3. **Docker Issues:**
   ```bash
   # On VM, check Docker
   docker ps
   docker-compose version
   ```

4. **Environment File Missing:**
   ```bash
   # Ensure .env.production exists on VM
   ls -la /opt/go-backend-api/.env.production
   ```

### CI Fails

1. **Test Failures:**
   - Check test logs in Actions tab
   - Run tests locally: `make test`

2. **Build Failures:**
   - Check Go version compatibility
   - Verify dependencies: `go mod tidy`

3. **Lint Failures:**
   - Run locally: `golangci-lint run`
   - Fix reported issues

## Security Best Practices

1. **Never commit secrets** - Always use GitHub Secrets
2. **Use SSH keys** - Don't use passwords for SSH
3. **Limit SSH access** - Use dedicated deploy user with minimal permissions
4. **Rotate keys regularly** - Update SSH keys periodically
5. **Use environment protection** - Require approvals for production deployments
6. **Monitor deployments** - Review deployment logs regularly

## Customization

### Modify Deployment Script

Edit `.github/workflows/deploy.yml` to customize:
- Deployment steps
- Pre/post deployment hooks
- Notification settings
- Rollback procedures

### Add Notifications

Add notification steps to workflows:
- Slack notifications
- Email alerts
- Discord webhooks
- Custom webhooks

Example:
```yaml
- name: Notify Slack
  uses: 8398a7/action-slack@v3
  with:
    status: ${{ job.status }}
    text: 'Deployment completed!'
  env:
    SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK }}
```

## Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Docker Build Push Action](https://github.com/docker/build-push-action)
- [SSH Agent Action](https://github.com/webfactory/ssh-agent)

