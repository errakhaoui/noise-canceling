# ClearVox Build Troubleshooting

## Common Build Issues

### Issue: `fatal error: 'rnnoise.h' file not found`

**Symptom:**
```
# github.com/errakhaoui/noise-canceling/noise_canceller
noise_canceller/noise_canceller.go:7:10: fatal error: 'rnnoise.h' file not found
    7 | #include <rnnoise.h>
      |          ^~~~~~~~~~~
```

**Cause:**
The C compiler can't find the RNNoise header files installed by Homebrew.

**Solution:**

Set CGO flags to point to Homebrew installation:

```bash
# Get Homebrew prefix (different on Intel vs Apple Silicon)
BREW_PREFIX=$(brew --prefix)

# Set environment variables
export CGO_CFLAGS="-I${BREW_PREFIX}/include"
export CGO_LDFLAGS="-L${BREW_PREFIX}/lib"
export PKG_CONFIG_PATH="${BREW_PREFIX}/lib/pkgconfig"

# Now build
go build -o clearvox-gui gui_main.go
```

**Quick fix:**

```bash
# One-liner
export CGO_CFLAGS="-I$(brew --prefix)/include" CGO_LDFLAGS="-L$(brew --prefix)/lib" && go build -o clearvox-gui gui_main.go
```

**Verify Homebrew installation:**

```bash
# Check Homebrew prefix
brew --prefix

# Check if rnnoise is installed
ls -la $(brew --prefix)/include/rnnoise.h
ls -la $(brew --prefix)/lib/librnnoise*

# If not found, install
brew install rnnoise portaudio
```

---

### Issue: `ld: library not found for -lrnnoise`

**Symptom:**
```
ld: library not found for -lrnnoise
clang: error: linker command failed with exit code 1
```

**Cause:**
The linker can't find the RNNoise library.

**Solution:**

Same as above - set `CGO_LDFLAGS`:

```bash
export CGO_LDFLAGS="-L$(brew --prefix)/lib"
go build -o clearvox-gui gui_main.go
```

---

### Issue: `ld: library not found for -lportaudio`

**Symptom:**
```
ld: library not found for -lportaudio
```

**Cause:**
PortAudio is not installed or not found by the linker.

**Solution:**

```bash
# Install PortAudio
brew install portaudio

# Set CGO flags
export CGO_LDFLAGS="-L$(brew --prefix)/lib"
go build -o clearvox-gui gui_main.go
```

---

### Issue: Different Homebrew paths on Intel vs Apple Silicon

**Homebrew paths:**
- **Intel Mac**: `/usr/local`
- **Apple Silicon**: `/opt/homebrew`

**Solution:**

Always use `brew --prefix` to get the correct path:

```bash
BREW_PREFIX=$(brew --prefix)
export CGO_CFLAGS="-I${BREW_PREFIX}/include"
export CGO_LDFLAGS="-L${BREW_PREFIX}/lib"
```

---

### Issue: GitHub Actions build fails

**Symptom:**
CI passes but release workflow fails with header not found error.

**Cause:**
The release workflow runs on macOS and needs CGO flags set.

**Solution:**

Already fixed in `.github/workflows/release.yml`:

```yaml
- name: Set up CGO flags for Homebrew libraries
  run: |
    BREW_PREFIX=$(brew --prefix)
    echo "CGO_CFLAGS=-I${BREW_PREFIX}/include" >> $GITHUB_ENV
    echo "CGO_LDFLAGS=-L${BREW_PREFIX}/lib" >> $GITHUB_ENV
```

---

## Build Script with CGO flags

Create a `build.sh` script:

```bash
#!/bin/bash
set -e

echo "=== Building ClearVox ==="

# Set up environment
BREW_PREFIX=$(brew --prefix)
export CGO_CFLAGS="-I${BREW_PREFIX}/include"
export CGO_LDFLAGS="-L${BREW_PREFIX}/lib"
export PKG_CONFIG_PATH="${BREW_PREFIX}/lib/pkgconfig"

echo "Homebrew prefix: $BREW_PREFIX"
echo "CGO_CFLAGS: $CGO_CFLAGS"
echo "CGO_LDFLAGS: $CGO_LDFLAGS"

# Build
go build -o clearvox-gui gui_main.go

echo "✓ Build successful"
echo "Binary: clearvox-gui"
```

Make it executable:

```bash
chmod +x build.sh
./build.sh
```

---

## Environment Setup for Development

Add to your `~/.zshrc` or `~/.bashrc`:

```bash
# ClearVox development environment
export CGO_CFLAGS="-I$(brew --prefix)/include"
export CGO_LDFLAGS="-L$(brew --prefix)/lib"
export PKG_CONFIG_PATH="$(brew --prefix)/lib/pkgconfig"
```

Then reload:

```bash
source ~/.zshrc  # or source ~/.bashrc
```

---

## Verifying Your Setup

Run this script to check everything:

```bash
#!/bin/bash
echo "=== ClearVox Build Environment Check ==="
echo ""

# Check Go
echo "Go version:"
go version
echo ""

# Check Homebrew
echo "Homebrew prefix:"
brew --prefix
echo ""

# Check RNNoise
echo "RNNoise installation:"
if [ -f "$(brew --prefix)/include/rnnoise.h" ]; then
    echo "✓ rnnoise.h found"
else
    echo "✗ rnnoise.h NOT found - run: brew install rnnoise"
fi
if [ -f "$(brew --prefix)/lib/librnnoise.0.dylib" ]; then
    echo "✓ librnnoise found"
else
    echo "✗ librnnoise NOT found"
fi
echo ""

# Check PortAudio
echo "PortAudio installation:"
if [ -f "$(brew --prefix)/lib/libportaudio.dylib" ]; then
    echo "✓ libportaudio found"
else
    echo "✗ libportaudio NOT found - run: brew install portaudio"
fi
echo ""

# Check CGO flags
echo "CGO environment:"
echo "CGO_CFLAGS: ${CGO_CFLAGS:-'not set'}"
echo "CGO_LDFLAGS: ${CGO_LDFLAGS:-'not set'}"
echo ""

echo "=== Recommended CGO flags ==="
echo "export CGO_CFLAGS=\"-I$(brew --prefix)/include\""
echo "export CGO_LDFLAGS=\"-L$(brew --prefix)/lib\""
```

---

## Quick Reference

**Install dependencies:**
```bash
brew install rnnoise portaudio
```

**Build with CGO flags:**
```bash
export CGO_CFLAGS="-I$(brew --prefix)/include"
export CGO_LDFLAGS="-L$(brew --prefix)/lib"
go build -o clearvox-gui gui_main.go
```

**Full build:**
```bash
rm -rf ClearVox.app clearvox-gui && \
export CGO_CFLAGS="-I$(brew --prefix)/include" && \
export CGO_LDFLAGS="-L$(brew --prefix)/lib" && \
go build -o clearvox-gui gui_main.go && \
~/bin/fyne package -os darwin --executable ./clearvox-gui && \
./create-dmg-installer.sh
```
