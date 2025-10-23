#!/bin/bash

set -e

echo "=== ClearVox App Installer ==="
echo ""

# Check if running on macOS
if [[ "$OSTYPE" != "darwin"* ]]; then
    echo "Error: This installer is for macOS only"
    exit 1
fi

# Check if Homebrew is installed
if ! command -v brew &> /dev/null; then
    echo "Homebrew not found. Installing Homebrew..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
else
    echo "✓ Homebrew is installed"
fi

# Check and install RNNoise
if brew list rnnoise &> /dev/null; then
    echo "✓ RNNoise is installed"
else
    echo "Installing RNNoise..."
    brew install rnnoise
fi

# Check and install PortAudio
if brew list portaudio &> /dev/null; then
    echo "✓ PortAudio is installed"
else
    echo "Installing PortAudio..."
    brew install portaudio
fi

# Check and install BlackHole (optional)
echo ""
echo "ClearVox can create a virtual microphone for video calls using BlackHole."
echo "Note: ClearVox will automatically offer to install BlackHole on first launch."
read -p "Install BlackHole now? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if brew list blackhole-2ch &> /dev/null; then
        echo "✓ BlackHole is already installed"
    else
        echo "Installing BlackHole..."
        brew install blackhole-2ch
        echo ""
        echo "✓ ClearVox Virtual Microphone ready!"
    fi
else
    echo "Skipping BlackHole installation. ClearVox will prompt you on first launch."
fi

echo ""
echo "=== Dependencies installed successfully ==="
echo ""
echo "Next steps:"
echo "1. Build the app: go build -o clearvox-gui gui_main.go"
echo "2. Package the app: ~/bin/fyne package -os darwin --executable ./clearvox-gui"
echo "3. Install to Applications: mv 'ClearVox.app' /Applications/"
echo "4. Launch: open '/Applications/ClearVox.app'"
echo ""
