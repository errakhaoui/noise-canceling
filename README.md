# ClearVox

Real-time noise cancellation for Go using RNNoise. Create a crystal-clear virtual microphone for your video calls. Available as GUI or CLI application.

## 📥 Download & Quick Start

### For End Users

**[⬇️ Download ClearVox-Installer.dmg](https://github.com/errakhaoui/noise-canceling/releases/latest)**

**Professional PKG Installation (30 seconds):**
1. Open the downloaded DMG file
2. Double-click **ClearVox-Installer.pkg**
3. Follow the installer prompts and enter your password
4. Launch ClearVox from Applications
5. Done! ✨

**✅ Professional installer** - Automated installation, no drag-and-drop needed. Everything is bundled.

### For Developers

See [Building from Source](#building-clearvox) below.

### Distribution Methods

ClearVox is distributed via **GitHub Releases** (industry standard for desktop apps).

For details on distribution options (Homebrew Cask, GitHub Packages, etc.), see [DISTRIBUTION.md](DISTRIBUTION.md).

---

## Requirements (For Developers)

- Go 1.22.3+
- RNNoise library
- PortAudio library

### Install Dependencies

**macOS:**
```bash
brew install rnnoise portaudio
# BlackHole is automatically installed by ClearVox on first launch
```

**Ubuntu/Debian:**
```bash
sudo apt-get install librnnoise-dev portaudio19-dev
```

**Arch Linux:**
```bash
sudo pacman -S rnnoise portaudio
```

## Installation

**Quick install with dependencies (macOS):**
```bash
git clone https://github.com/errakhaoui/noise-canceling.git
cd noise-canceling
./install.sh
go mod download
```

**Manual installation:**
```bash
git clone https://github.com/errakhaoui/noise-canceling.git
cd noise-canceling
go mod download
```
Then install dependencies manually (see requirements above).

## Build & Run

**GUI (recommended):**
```bash
go build -o clearvox-gui gui_main.go
./clearvox-gui
```

**CLI:**
```bash
go build -o clearvox example.go
./clearvox
```

## Build macOS App Bundle

### For End Users (Non-Technical)

**✨ Professional PKG Installer - One-Click Installation!**

ClearVox uses a professional macOS installer package (.pkg) for seamless, automated installation.

**Installation via DMG:**
1. Download `ClearVox-Installer.dmg` from [GitHub Releases](https://github.com/errakhaoui/noise-canceling/releases)
2. Open the DMG file
3. Read the README.txt (optional)
4. Double-click **"ClearVox-Installer.pkg"**
5. Follow the macOS Installer prompts
6. Enter your password when requested
7. Click "Install"
8. Done! ClearVox is now in your Applications folder

**First Launch:**
1. Open ClearVox from Applications or Spotlight
2. Grant microphone permission when prompted (automatic macOS dialog)
3. (Optional) Click "Yes" when asked to install virtual microphone driver

**What's Bundled in the PKG:**
- ✅ RNNoise library (noise cancellation engine)
- ✅ PortAudio library (audio I/O)
- ✅ All required dependencies
- ✅ Professional installer with welcome & completion screens
- ✅ Automatic installation to /Applications
- ✅ Proper permissions setup

**Benefits of PKG Installer:**
- ✅ One-click automated installation
- ✅ Professional macOS installer UI
- ✅ No drag-and-drop needed
- ✅ Handles permissions automatically
- ✅ Standard installation process
- ✅ Works on any Mac (macOS 10.11+)

**Optional Virtual Microphone:**
- BlackHole is automatically downloaded and installed when you first launch ClearVox
- Only needed if you want to use ClearVox as a virtual microphone in video calls
- Requires clicking "Yes" and entering your password once

### For Developers

#### Prerequisites

1. **Install Go 1.22.3+**
   ```bash
   brew install go
   ```

2. **Install system dependencies** (RNNoise, PortAudio)
   ```bash
   ./install.sh
   ```
   This will install:
   - RNNoise library for noise cancellation
   - PortAudio library for audio I/O
   - Optionally: BlackHole 2ch for virtual microphone

3. **Install Fyne packaging tool**
   ```bash
   go install fyne.io/tools/cmd/fyne@latest
   ```
   This installs to `~/bin/fyne`

#### Building ClearVox

**Step 1: Build the executable**
```bash
go build -o clearvox-gui gui_main.go
```
This creates a `clearvox-gui` binary (~30MB)

**Step 2: Package as macOS app bundle**
```bash
~/bin/fyne package -os darwin --executable ./clearvox-gui
```
This creates `ClearVox.app` with proper macOS app structure

**Step 3: Create distributable DMG with PKG installer**
```bash
./create-dmg-installer.sh
```

This script automatically:
- ✅ Adds microphone permissions to `Info.plist`
- ✅ Bundles RNNoise and PortAudio libraries into the app
- ✅ Fixes library paths to use bundled versions
- ✅ Re-signs the app with ad-hoc signature
- ✅ Creates professional PKG installer with welcome/conclusion screens
- ✅ Packages PKG inside DMG with README
- ✅ Creates `ClearVox-Installer.dmg` (~20MB)

**One-line build command:**
```bash
rm -rf ClearVox.app clearvox-gui && \
go build -o clearvox-gui gui_main.go && \
~/bin/fyne package -os darwin --executable ./clearvox-gui && \
./create-dmg-installer.sh
```

#### Testing the Build

**Option 1: Test app directly**
```bash
open ClearVox.app
```

**Option 2: Test DMG installer**
```bash
open ClearVox-Installer.dmg
# Drag ClearVox.app to Applications and launch
```

**Testing checklist:**
- ✅ App launches without crashing
- ✅ macOS prompts for microphone permission
- ✅ Audio devices are listed in dropdowns
- ✅ App offers to install BlackHole (if not installed)
- ✅ Audio processing works when Start is clicked

#### Code Signing for Distribution

**For local testing:**
- Script uses ad-hoc signature (`-`) - works on your Mac

**For public distribution:**
1. Get Apple Developer certificate ($99/year)
2. Edit `create-dmg-installer.sh` line 112
3. Replace:
   ```bash
   codesign --force --deep --sign - "$APP_BUNDLE"
   ```
   With:
   ```bash
   codesign --force --deep --sign "Developer ID Application: Your Name (TEAMID)" "$APP_BUNDLE"
   ```
4. Optionally notarize with Apple:
   ```bash
   xcrun notarytool submit ClearVox-Installer.dmg \
     --apple-id your@email.com \
     --password app-specific-password \
     --team-id TEAMID
   ```

#### Build Artifacts

After successful build:
```
clearvox-gui              # Executable binary (30MB)
ClearVox.app/             # macOS app bundle
ClearVox-Installer.pkg    # Professional PKG installer (30MB)
ClearVox-Installer.dmg    # DMG containing PKG + README (20MB)
```

**DMG Contents:**
```
ClearVox-Installer.dmg/
├── ClearVox-Installer.pkg    # Double-click to install
└── README.txt                # Installation instructions
```

#### Publishing a Release (Fully Automated)

Releases are **automatically built and published** when you update the version in `FyneApp.toml` and merge to main.

**🚀 Super Simple Release Process:**

```bash
# 1. Update version in FyneApp.toml
# Edit the Version field: Version = "1.0.1"

# 2. Commit and push to main
git add FyneApp.toml
git commit -m "chore: bump version to 1.0.1"
git push origin main
```

**That's it!** 🎉 GitHub Actions automatically:
- ✅ Detects the version change in `FyneApp.toml`
- ✅ Builds the executable
- ✅ Packages as macOS app
- ✅ Creates DMG installer
- ✅ Signs the app
- ✅ Creates git tag (e.g., `v1.0.1`)
- ✅ Creates GitHub Release
- ✅ Uploads `ClearVox-Installer.dmg`

**How it works:**
- Workflow triggers **only** when `FyneApp.toml` is changed on `main` branch
- Version is extracted from `FyneApp.toml`
- Git tag is created automatically
- Release appears at: https://github.com/errakhaoui/noise-canceling/releases/latest

**Monitor the build:**
- Go to: https://github.com/errakhaoui/noise-canceling/actions
- Watch the "Release" workflow complete (~5 minutes)

**Example workflow:**

```bash
# Working on a feature
git checkout -b feat/new-feature
# ... make changes ...
git commit -m "feat: add new feature"
git push origin feat/new-feature

# Create PR and merge to main
# (No release triggered yet)

# When ready to release:
git checkout main
git pull
# Edit FyneApp.toml: Version = "1.1.0"
git add FyneApp.toml
git commit -m "chore: release v1.1.0"
git push origin main
# 🎉 Release automatically triggered!
```

**Manual trigger (optional):**

You can also trigger a release manually from GitHub:
1. Go to: https://github.com/errakhaoui/noise-canceling/actions/workflows/release.yml
2. Click "Run workflow"
3. Choose `main` branch and click "Run workflow"

**Release Checklist:**
- ✅ All tests passing on main branch
- ✅ Version number updated in `FyneApp.toml`
- ✅ Changes committed and pushed to main
- ✅ CI auto-generates release notes and creates release

#### Troubleshooting

**"ClearVox.app not found"**
- Run Step 2 first: `~/bin/fyne package -os darwin --executable ./clearvox-gui`

**"Code signature invalid"**
- The script automatically re-signs. If issues persist:
  ```bash
  codesign --force --deep --sign - ClearVox.app
  ```

**"Microphone permission not requested"**
- Reset permissions:
  ```bash
  tccutil reset Microphone com.errakhaoui.clearvox
  ```
- Rebuild the DMG to ensure `NSMicrophoneUsageDescription` is present

**Library not found errors**
- Ensure dependencies are installed: `./install.sh`
- Check library paths:
  ```bash
  otool -L clearvox-gui
  ```

#### CI/CD Workflows

ClearVox uses GitHub Actions for automated testing and releases:

**`.github/workflows/ci.yml` - Continuous Integration**
- Triggers on: Pull requests
- Runs: Tests, linting, code formatting checks
- Platform: Ubuntu (for testing)

**`.github/workflows/release.yml` - Automated Releases**
- Triggers on: Push to `main` when `FyneApp.toml` changes, or manual dispatch
- Runs: Build → Package → DMG creation → Tag creation → GitHub Release
- Platform: macOS (required for DMG creation)
- Output: `ClearVox-Installer.dmg` uploaded to GitHub Releases
- Smart: Only creates release if version changed (prevents duplicate releases)

**Workflow Status:**

[![CI](https://github.com/errakhaoui/noise-canceling/actions/workflows/ci.yml/badge.svg)](https://github.com/errakhaoui/noise-canceling/actions/workflows/ci.yml)
[![Release](https://github.com/errakhaoui/noise-canceling/actions/workflows/release.yml/badge.svg)](https://github.com/errakhaoui/noise-canceling/actions/workflows/release.yml)

---

## Virtual Microphone Setup (Google Meet/Zoom)

ClearVox creates a crystal-clear virtual microphone for your video calls with **zero manual setup**.

**Automatic Setup (First Launch):**
1. Launch ClearVox
2. When prompted, click "Yes" to install BlackHole audio driver
3. Enter your admin password when macOS requests it
4. Restart ClearVox after installation completes

**Using Your Virtual Microphone:**
1. In ClearVox:
   - Input Device: Your physical microphone
   - Output Device: BlackHole 2ch (your virtual microphone)
   - Monitor Device: Your headphones (optional, to hear yourself)
   - Click Start
2. In Google Meet/Zoom: Set microphone to "BlackHole 2ch"

Your audio now flows: **Physical Mic → ClearVox (noise canceling) → Virtual Mic → Meeting**

**Manual Installation (Optional):**
If you prefer to install BlackHole manually: `brew install blackhole-2ch`

**Note:** You can use any virtual audio device (Loopback, Soundflower, etc.) as your ClearVox Virtual Microphone output.

## CLI Usage

```bash
# List devices
./clearvox -list-devices

# Use default output
./clearvox

# Route to virtual microphone with monitoring
./clearvox -device blackhole -monitor-device headphones

# Toggle noise cancellation: type 't' + Enter
```

## Testing

```bash
# Run tests
go test ./... -short -v

# With coverage
go test ./... -short -v -race -coverprofile=coverage.out
```

## Project Structure

```
clearvox/
├── gui_main.go              # GUI entry point
├── example.go               # CLI entry point
├── gui/                     # GUI components
├── input/                   # Microphone capture
├── noise_canceller/         # RNNoise integration
└── output/                  # Audio playback
```

## License

MIT
