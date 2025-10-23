# ClearVox Distribution Guide

## Current Distribution: GitHub Releases ✅ (Recommended)

ClearVox is currently distributed via **GitHub Releases**, which is the industry-standard approach for macOS desktop applications.

### Why GitHub Releases?

- ✅ **Industry Standard** - Used by VS Code, Docker Desktop, etc.
- ✅ **User-Friendly** - Direct download links, easy to find
- ✅ **Automatic** - CI/CD creates releases automatically
- ✅ **Professional** - Release notes, version history, assets
- ✅ **No Setup** - Users just download and install
- ✅ **Free** - Unlimited public releases

### Current Workflow

```
User downloads → ClearVox-Installer.dmg
Opens DMG → Double-clicks ClearVox-Installer.pkg
macOS Installer runs → Installed to /Applications
```

**Download URL:** https://github.com/errakhaoui/noise-canceling/releases/latest

---

## Alternative: Homebrew Cask (Coming Soon)

For CLI users who prefer `brew install`, we can publish a Homebrew Cask.

### What is Homebrew Cask?

Homebrew Cask extends Homebrew to install GUI macOS applications via command line.

### Installation with Homebrew

Once published:
```bash
# Install ClearVox
brew install --cask clearvox

# Update ClearVox
brew upgrade --cask clearvox

# Uninstall ClearVox
brew uninstall --cask clearvox
```

### How to Publish to Homebrew

1. **First, create a release on GitHub** (already automated)
2. **Submit Cask to Homebrew** (one-time process)

**Steps:**

```bash
# 1. Fork homebrew-cask
# Visit: https://github.com/Homebrew/homebrew-cask
# Click "Fork"

# 2. Clone your fork
git clone https://github.com/YOUR-USERNAME/homebrew-cask.git
cd homebrew-cask

# 3. Create cask file
cp ~/noise-canceling/Casks/clearvox.rb Casks/c/clearvox.rb

# 4. Update SHA256 (after first release)
# Download your DMG and run:
shasum -a 256 ClearVox-Installer.dmg
# Update the sha256 value in clearvox.rb

# 5. Test the cask locally
brew install --cask Casks/c/clearvox.rb

# 6. Create PR
git checkout -b clearvox
git add Casks/c/clearvox.rb
git commit -m "clearvox: add new cask"
git push origin clearvox

# 7. Create Pull Request on GitHub
# Visit: https://github.com/Homebrew/homebrew-cask/compare
# Create PR from your fork
```

**Cask Requirements:**
- App must be released on GitHub Releases
- SHA256 checksum must be provided
- Cask must follow Homebrew style guidelines
- Open source or freeware

### Cask File Location

The cask definition is in: `Casks/clearvox.rb`

To submit to Homebrew Cask: https://github.com/Homebrew/homebrew-cask

---

## GitHub Packages (Not Recommended for Desktop Apps)

GitHub Packages is designed for:
- Docker containers
- npm/Maven/NuGet packages
- Programmatic dependency management

**Why NOT use GitHub Packages for ClearVox:**
- ❌ Not designed for desktop app distribution
- ❌ Requires authentication for downloads
- ❌ No user-friendly download page
- ❌ Adds unnecessary complexity
- ❌ Users expect GitHub Releases for apps

**GitHub Packages is great for:**
- Libraries and SDKs
- Container images
- Language-specific packages

**GitHub Releases is perfect for:**
- Desktop applications (like ClearVox)
- DMG/PKG installers
- User-facing downloads

---

## Distribution Comparison

| Method | User Experience | Setup Complexity | Updates | Best For |
|--------|----------------|------------------|---------|----------|
| **GitHub Releases** | ⭐⭐⭐⭐⭐ Download DMG | ⭐⭐⭐⭐⭐ Automatic | Manual | **Most users** |
| **Homebrew Cask** | ⭐⭐⭐⭐ `brew install` | ⭐⭐⭐ Submit once | `brew upgrade` | **CLI users** |
| **GitHub Packages** | ⭐⭐ Complex auth | ⭐⭐ Complex setup | Programmatic | **Not suitable** |

---

## Recommended Approach

### Phase 1: GitHub Releases (Current) ✅

**Status:** Implemented and automated via CI/CD

**Users:**
1. Visit https://github.com/errakhaoui/noise-canceling/releases
2. Download ClearVox-Installer.dmg
3. Install

**You:**
1. Update version in FyneApp.toml
2. Push to main
3. Release created automatically

### Phase 2: Add Homebrew Cask (Optional)

**For advanced users who want CLI installation**

**Steps:**
1. Wait for first stable release (v1.0.0)
2. Get SHA256 of DMG
3. Update `Casks/clearvox.rb` with SHA256
4. Submit PR to homebrew-cask
5. Once merged, users can `brew install --cask clearvox`

**Benefits:**
- CLI users happy
- Automatic updates via Homebrew
- Increased visibility
- Professional image

**Timeline:**
- After first stable release
- ~1 week for Homebrew review/approval
- Then available to all Homebrew users

---

## Current Status

✅ **GitHub Releases** - Fully automated, ready for users
⏳ **Homebrew Cask** - Ready to submit after v1.0.0 release
❌ **GitHub Packages** - Not applicable for desktop apps

---

## For End Users

**Preferred Method:** GitHub Releases
- Download: https://github.com/errakhaoui/noise-canceling/releases/latest
- One-click DMG installation
- No terminal required

**Coming Soon:** Homebrew Cask
- Command: `brew install --cask clearvox`
- For CLI enthusiasts

---

## Summary

**Keep using GitHub Releases** - It's the right choice for ClearVox. It's:
- Industry standard for desktop apps
- User-friendly
- Fully automated
- Professional

**Optionally add Homebrew Cask** - For CLI users who want `brew install`.

**Skip GitHub Packages** - Not designed for desktop application distribution.
