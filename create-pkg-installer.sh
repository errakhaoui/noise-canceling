#!/bin/bash

set -e

echo "=== Creating ClearVox PKG Installer ==="
echo ""

APP_NAME="ClearVox"
APP_BUNDLE="${APP_NAME}.app"
PKG_NAME="ClearVox-Installer"
VERSION="1.0.0"
IDENTIFIER="com.errakhaoui.clearvox"

# Check if app bundle exists
if [ ! -d "$APP_BUNDLE" ]; then
    echo "Error: $APP_BUNDLE not found. Please build it first."
    exit 1
fi

# Create temporary directories
TEMP_DIR=$(mktemp -d)
BUILD_DIR="$TEMP_DIR/build"
SCRIPTS_DIR="$TEMP_DIR/scripts"
RESOURCES_DIR="$TEMP_DIR/resources"

mkdir -p "$BUILD_DIR"
mkdir -p "$SCRIPTS_DIR"
mkdir -p "$RESOURCES_DIR"

echo "Step 1: Preparing app bundle for packaging..."

# Copy app to temporary location
cp -R "$APP_BUNDLE" "$BUILD_DIR/"

echo "Step 2: Creating post-install script..."

# Create post-install script to fix permissions and open app
cat > "$SCRIPTS_DIR/postinstall" << 'EOF'
#!/bin/bash

# Fix permissions
chown -R root:wheel "/Applications/ClearVox.app"
chmod -R 755 "/Applications/ClearVox.app"

# Open the app after installation (optional - user will see welcome dialog)
# Uncomment to auto-open:
# open "/Applications/ClearVox.app"

exit 0
EOF

chmod +x "$SCRIPTS_DIR/postinstall"

echo "Step 3: Creating welcome and conclusion messages..."

# Create welcome message
cat > "$RESOURCES_DIR/welcome.html" << 'EOF'
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, sans-serif; margin: 20px; }
        h1 { color: #1d1d1f; font-size: 32px; font-weight: 600; }
        p { color: #1d1d1f; font-size: 14px; line-height: 1.6; }
        .highlight { background: #f5f5f7; padding: 15px; border-radius: 8px; margin: 15px 0; }
        ul { margin: 10px 0; padding-left: 20px; }
        li { margin: 8px 0; }
    </style>
</head>
<body>
    <h1>Welcome to ClearVox</h1>
    <p>ClearVox provides real-time noise cancellation for crystal-clear audio in your video calls.</p>

    <div class="highlight">
        <strong>✨ What's Included:</strong>
        <ul>
            <li>🎤 Real-time noise cancellation</li>
            <li>🖥️ Beautiful native macOS interface</li>
            <li>🎧 Virtual microphone for video calls</li>
            <li>🚀 Auto-install BlackHole audio driver</li>
            <li>📦 All dependencies bundled</li>
        </ul>
    </div>

    <p><strong>Installation Size:</strong> Approximately 32 MB</p>
    <p><strong>Installation Location:</strong> /Applications/ClearVox.app</p>
</body>
</html>
EOF

# Create conclusion message
cat > "$RESOURCES_DIR/conclusion.html" << 'EOF'
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, sans-serif; margin: 20px; }
        h1 { color: #1d1d1f; font-size: 32px; font-weight: 600; }
        p { color: #1d1d1f; font-size: 14px; line-height: 1.6; }
        .steps { background: #f5f5f7; padding: 15px; border-radius: 8px; margin: 15px 0; }
        ol { margin: 10px 0; padding-left: 25px; }
        li { margin: 10px 0; }
    </style>
</head>
<body>
    <h1>Installation Complete! 🎉</h1>
    <p>ClearVox has been successfully installed to your Applications folder.</p>

    <div class="steps">
        <strong>Next Steps:</strong>
        <ol>
            <li>Launch <strong>ClearVox</strong> from Applications or Spotlight</li>
            <li>Grant <strong>microphone permission</strong> when prompted</li>
            <li>Click <strong>"Yes"</strong> to install virtual microphone driver (optional)</li>
            <li>Configure your audio devices and click <strong>Start</strong></li>
        </ol>
    </div>

    <p><strong>Need Help?</strong></p>
    <p>Visit: <a href="https://github.com/errakhaoui/noise-canceling">github.com/errakhaoui/noise-canceling</a></p>
</body>
</html>
EOF

echo "Step 4: Building component package..."

# Build the component package
pkgbuild --root "$BUILD_DIR" \
         --install-location "/Applications" \
         --scripts "$SCRIPTS_DIR" \
         --identifier "$IDENTIFIER" \
         --version "$VERSION" \
         "$TEMP_DIR/ClearVox-component.pkg"

echo "Step 5: Creating distribution XML..."

# Create distribution definition for product
cat > "$TEMP_DIR/Distribution.xml" << EOF
<?xml version="1.0" encoding="utf-8"?>
<installer-gui-script minSpecVersion="1">
    <title>ClearVox</title>
    <organization>com.errakhaoui</organization>
    <domains enable_localSystem="true"/>
    <options customize="never" require-scripts="false" hostArchitectures="arm64,x86_64"/>

    <!-- Define documents displayed at various steps -->
    <welcome file="welcome.html" mime-type="text/html"/>
    <conclusion file="conclusion.html" mime-type="text/html"/>

    <!-- License (optional) -->
    <!-- <license file="license.txt" mime-type="text/plain"/> -->

    <!-- Background image (optional) -->
    <!-- <background file="background.png" mime-type="image/png" alignment="bottomleft" scaling="proportional"/> -->

    <!-- List all component packages -->
    <pkg-ref id="$IDENTIFIER"/>

    <!-- Define the order and visibility of components -->
    <options customize="never" require-scripts="false"/>
    <choices-outline>
        <line choice="default">
            <line choice="$IDENTIFIER"/>
        </line>
    </choices-outline>

    <choice id="default"/>
    <choice id="$IDENTIFIER" visible="false">
        <pkg-ref id="$IDENTIFIER"/>
    </choice>

    <pkg-ref id="$IDENTIFIER" version="$VERSION" onConclusion="none">ClearVox-component.pkg</pkg-ref>
</installer-gui-script>
EOF

echo "Step 6: Building final product package..."

# Build the final product package
productbuild --distribution "$TEMP_DIR/Distribution.xml" \
             --resources "$RESOURCES_DIR" \
             --package-path "$TEMP_DIR" \
             "${PKG_NAME}.pkg"

echo "Step 7: Signing package (optional)..."

# Sign the package (if you have a Developer ID Installer certificate)
# Uncomment and replace with your certificate name:
# productsign --sign "Developer ID Installer: Your Name (TEAMID)" \
#             "${PKG_NAME}.pkg" \
#             "${PKG_NAME}-signed.pkg"
# mv "${PKG_NAME}-signed.pkg" "${PKG_NAME}.pkg"

echo "Step 8: Cleaning up..."

rm -rf "$TEMP_DIR"

echo ""
echo "=== PKG Installer created successfully ==="
echo "File: ${PKG_NAME}.pkg"
echo ""
echo "To install: double-click ${PKG_NAME}.pkg"
echo ""
