# Quick Start: GitHub Actions CI/CD Setup

## 1. Generate SSH Key (One-time setup)

```bash
ssh-keygen -t ed25519 -C "github-actions" -f ~/.ssh/github_actions_deploy
ssh-copy-id -i ~/.ssh/github_actions_deploy.pub user@your-vm-host
cat ~/.ssh/github_actions_deploy  # Copy this output
```

## 2. Add GitHub Secrets

Go to: **Repository → Settings → Secrets and variables → Actions**

Add these secrets:

- `VM_HOST` = `your-vm-ip-or-domain`
- `VM_USER` = `your-ssh-username`
- `VM_SSH_PRIVATE_KEY` = (paste the private key from step 1)

## 3. Prepare VM

```bash
# On your VM
sudo mkdir -p /opt/go-backend-api
sudo chown $USER:$USER /opt/go-backend-api
cd /opt/go-backend-api
cp env.production.example .env.production
nano .env.production  # Edit with your values
```

## 4. Test It!

```bash
# Push to main branch
git push origin main

# Or trigger manually:
# Go to Actions → Deploy → Run workflow
```

That's it! 🎉

For detailed documentation, see [README.md](./README.md)

