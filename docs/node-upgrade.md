# Node.js Upgrade Guide

This project now requires **Node.js v24.22.1 (LTS)** or later.

## Quick Upgrade

### Using nvm (Recommended)

```bash
# Install/use Node.js 24
nvm install 24
nvm use 24

# Verify version
node --version  # Should show v24.22.1 or later

# Set as default (optional)
nvm alias default 24
```

The project includes `.nvmrc` files, so you can simply:
```bash
# In project root or frontend/
nvm use
```

### Using Official Installer

1. Download Node.js 24 LTS from https://nodejs.org/
2. Install for your platform
3. Verify: `node --version`

### Using Package Managers

**macOS (Homebrew)**:
```bash
brew install node@24
brew link node@24
```

**Ubuntu/Debian**:
```bash
curl -fsSL https://deb.nodesource.com/setup_24.x | sudo -E bash -
sudo apt-get install -y nodejs
```

**Windows (Chocolatey)**:
```bash
choco install nodejs-lts
```

## After Upgrading

1. **Reinstall dependencies**:
   ```bash
   cd frontend
   rm -rf node_modules package-lock.json
   npm install
   ```

2. **Verify everything works**:
   ```bash
   ./run.sh init
   ./run.sh dev
   ```

## What Changed

- ✅ **Dockerfile**: Updated to `node:24-alpine`
- ✅ **package.json**: Added `engines` field requiring Node ≥24.22.1
- ✅ **.nvmrc**: Created for nvm users
- ✅ **run.sh**: Now checks for Node 24+
- ✅ **Documentation**: Updated all references

## Why Upgrade?

Node.js 24 is the latest LTS release with:
- 🚀 **Better performance**: V8 engine improvements
- 🔒 **Enhanced security**: Latest security patches
- ✨ **New features**: Modern JavaScript support
- 📦 **npm 10**: Latest package manager
- 🛠️ **Long-term support**: Maintained until 2027

## Troubleshooting

### "node: command not found"

Node.js is not installed or not in PATH.
- Install from https://nodejs.org/
- Or use nvm: `curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash`

### Version check fails in run.sh

```bash
# Check your current version
node --version

# If less than v24, upgrade using one of the methods above
```

### npm version mismatch

```bash
# Update npm to latest
npm install -g npm@latest

# Verify version
npm --version  # Should be 10.0.0+
```

### Permission errors on macOS/Linux

```bash
# Fix npm global permissions
mkdir ~/.npm-global
npm config set prefix '~/.npm-global'
echo 'export PATH=~/.npm-global/bin:$PATH' >> ~/.bashrc
source ~/.bashrc
```

## Using nvm (Recommended for Developers)

If you work on multiple Node.js projects with different versions:

1. **Install nvm**:
   ```bash
   curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
   ```

2. **Install Node 24**:
   ```bash
   nvm install 24
   ```

3. **Auto-switch with .nvmrc**:
   The project includes `.nvmrc` files that automatically switch to the correct version:
   ```bash
   cd /path/to/ovnkube-viewer
   nvm use  # Automatically uses v24.22.1
   ```

4. **Add to shell profile** (optional):
   ```bash
   # Add to ~/.bashrc or ~/.zshrc
   autoload -U add-zsh-hook
   load-nvmrc() {
     if [[ -f .nvmrc && -r .nvmrc ]]; then
       nvm use
     fi
   }
   add-zsh-hook chpwd load-nvmrc
   ```

## Verification

After upgrading, verify your setup:

```bash
# Check versions
node --version    # v24.22.1 or later
npm --version     # 10.0.0 or later

# Check run.sh
./run.sh init     # Should pass dependency checks

# Run the app
./run.sh dev      # Should start without errors
```

## Docker Users

If you only use Docker, you don't need to upgrade your local Node.js:
- The Dockerfile uses `node:24-alpine`
- Container has the correct version automatically

```bash
./run.sh docker   # Uses Node 24 in container
```

## References

- [Node.js Official Site](https://nodejs.org/)
- [Node.js 24 Release Notes](https://nodejs.org/en/blog/release/v24.0.0)
- [nvm (Node Version Manager)](https://github.com/nvm-sh/nvm)
- [npm Documentation](https://docs.npmjs.com/)

---

**Updated**: October 26, 2026  
**Required Version**: Node.js ≥ v24.22.1 (LTS)
