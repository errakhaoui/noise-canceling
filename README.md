# noise-canceling

Real-time noise cancellation for Go using RNNoise. Available as GUI or CLI application.

## Requirements

- Go 1.22.3+
- RNNoise library
- PortAudio library

### Install Dependencies

**macOS:**
```bash
brew install rnnoise portaudio
brew install blackhole-2ch  # Optional: for virtual microphone
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

```bash
git clone https://github.com/errakhaoui/noise-canceling.git
cd noise-canceling
go mod download
```

## Build & Run

**GUI (recommended):**
```bash
go build -o noise-canceling-gui gui_main.go
./noise-canceling-gui
```

**CLI:**
```bash
go build -o noise-canceling example.go
./noise-canceling
```

## Virtual Microphone Setup (Google Meet/Zoom)

1. Install BlackHole: `brew install blackhole-2ch`
2. Launch GUI application
3. Configure:
   - Input Device: Your physical microphone
   - Output Device: BlackHole 2ch
   - Monitor Device: Your headphones (optional)
4. Click Start
5. In Google Meet/Zoom: Set microphone to "BlackHole 2ch"

## CLI Usage

```bash
# List devices
./noise-canceling -list-devices

# Use default output
./noise-canceling

# Virtual microphone with monitoring
./noise-canceling -device blackhole -monitor-device headphones

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
noise-canceling/
├── gui_main.go              # GUI entry point
├── example.go               # CLI entry point
├── gui/                     # GUI components
├── input/                   # Microphone capture
├── noise_canceller/         # RNNoise integration
└── output/                  # Audio playback
```

## License

MIT
