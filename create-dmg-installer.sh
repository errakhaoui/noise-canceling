#!/bin/bash

set -e

echo "=== Creating macOS DMG Installer ==="
echo ""

APP_NAME="ClearVox"
APP_BUNDLE="${APP_NAME}.app"
DMG_NAME="ClearVox-Installer"
VERSION="1.0.0"

# Check if app bundle exists
if [ ! -d "$APP_BUNDLE" ]; then
    echo "Error: $APP_BUNDLE not found. Please build it first."
    exit 1
fi

echo "Step 1: Adding microphone permissions..."

INFO_PLIST="$APP_BUNDLE/Contents/Info.plist"

# Check if microphone permission is already present
if ! grep -q "NSMicrophoneUsageDescription" "$INFO_PLIST"; then
    # Add microphone permission using PlistBuddy
    /usr/libexec/PlistBuddy -c "Add :NSMicrophoneUsageDescription string 'ClearVox needs access to your microphone to provide real-time noise cancellation for crystal-clear audio in video calls.'" "$INFO_PLIST" 2>/dev/null || true
    echo "  ✓ Microphone permission added"
else
    echo "  ✓ Microphone permission already present"
fi

echo "Step 2: Bundling dependencies into app..."

# Create Frameworks directory inside the app bundle
FRAMEWORKS_DIR="$APP_BUNDLE/Contents/Frameworks"
mkdir -p "$FRAMEWORKS_DIR"

# Find and copy RNNoise library
RNNOISE_PATHS=(
    "/usr/local/lib/librnnoise.0.dylib"
    "$(brew --prefix 2>/dev/null)/lib/librnnoise.0.dylib"
)

RNNOISE_FOUND=false
for path in "${RNNOISE_PATHS[@]}"; do
    if [ -f "$path" ]; then
        echo "  Copying RNNoise library from $path..."
        cp "$path" "$FRAMEWORKS_DIR/"
        RNNOISE_FOUND=true
        break
    fi
done

if [ "$RNNOISE_FOUND" = false ]; then
    echo "  Error: RNNoise library not found!"
    exit 1
fi

# Find and copy PortAudio library
PORTAUDIO_PATHS=(
    "/opt/homebrew/opt/portaudio/lib/libportaudio.2.dylib"
    "/usr/local/opt/portaudio/lib/libportaudio.2.dylib"
    "$(brew --prefix portaudio 2>/dev/null)/lib/libportaudio.2.dylib"
)

PORTAUDIO_FOUND=false
for path in "${PORTAUDIO_PATHS[@]}"; do
    if [ -f "$path" ]; then
        echo "  Copying PortAudio library from $path..."
        cp "$path" "$FRAMEWORKS_DIR/"
        PORTAUDIO_FOUND=true
        break
    fi
done

if [ "$PORTAUDIO_FOUND" = false ]; then
    echo "  Error: PortAudio library not found!"
    exit 1
fi

echo "Step 3: Fixing library paths..."

EXECUTABLE="$APP_BUNDLE/Contents/MacOS/clearvox-gui"

# Fix RNNoise path
install_name_tool -change "/usr/local/lib/librnnoise.0.dylib" \
    "@executable_path/../Frameworks/librnnoise.0.dylib" "$EXECUTABLE" 2>/dev/null || \
install_name_tool -change "$(brew --prefix)/lib/librnnoise.0.dylib" \
    "@executable_path/../Frameworks/librnnoise.0.dylib" "$EXECUTABLE" 2>/dev/null || true

# Fix PortAudio path
install_name_tool -change "$(brew --prefix portaudio)/lib/libportaudio.2.dylib" \
    "@executable_path/../Frameworks/libportaudio.2.dylib" "$EXECUTABLE" 2>/dev/null || \
install_name_tool -change "/usr/local/lib/libportaudio.2.dylib" \
    "@executable_path/../Frameworks/libportaudio.2.dylib" "$EXECUTABLE" 2>/dev/null || true

# Update library IDs
if [ -f "$FRAMEWORKS_DIR/librnnoise.0.dylib" ]; then
    install_name_tool -id "@executable_path/../Frameworks/librnnoise.0.dylib" \
        "$FRAMEWORKS_DIR/librnnoise.0.dylib"
fi

if [ -f "$FRAMEWORKS_DIR/libportaudio.2.dylib" ]; then
    install_name_tool -id "@executable_path/../Frameworks/libportaudio.2.dylib" \
        "$FRAMEWORKS_DIR/libportaudio.2.dylib"
fi

echo "Step 4: Re-signing application..."

# Remove existing signature (if any) and re-sign with ad-hoc signature
# For distribution, replace "-" with your Apple Developer certificate ID
codesign --force --deep --sign - "$APP_BUNDLE"

if [ $? -eq 0 ]; then
    echo "  ✓ Application signed successfully"
else
    echo "  ⚠ Warning: Signing failed. The app may not run on some systems."
fi

echo "Step 5: Creating DMG installer..."

# Remove old DMG if exists
rm -f "${DMG_NAME}.dmg"
rm -rf dmg_temp

# Create temporary directory for DMG contents
mkdir -p dmg_temp

# Copy app bundle to temp directory
cp -R "$APP_BUNDLE" dmg_temp/

# Create Applications symlink for drag-and-drop
ln -s /Applications dmg_temp/Applications

# Create a simple README
cat > dmg_temp/README.txt << EOF
ClearVox - Installation Instructions
=====================================

✨ ZERO SETUP REQUIRED - Everything is bundled!

Installation:
1. Drag "ClearVox.app" to the Applications folder
2. Open "ClearVox.app" from Applications or Spotlight
3. Grant microphone permission when prompted (automatic macOS dialog)
4. (Optional) Click "Yes" to install virtual microphone for video calls
5. Enter your admin password if installing virtual microphone

That's it! No Homebrew, no terminal, no dependencies to install.

What's Included:
✓ RNNoise library (noise cancellation)
✓ PortAudio library (audio I/O)
✓ All dependencies built-in

Works on any Mac running macOS 10.11 or newer.

For more information, visit:
https://github.com/errakhaoui/noise-canceling

Version: ${VERSION}
EOF

# Create DMG
echo "  Creating disk image..."
hdiutil create -volname "${APP_NAME}" \
    -srcfolder dmg_temp \
    -ov -format UDZO \
    "${DMG_NAME}.dmg"

# Cleanup
rm -rf dmg_temp

echo ""
echo "=== DMG Installer created successfully ==="
echo "File: ${DMG_NAME}.dmg"
echo ""
echo "The app now includes all required libraries (RNNoise, PortAudio)."
echo "Users can simply drag the app to Applications - no technical setup needed!"
echo ""
echo "To test: open ${DMG_NAME}.dmg"
