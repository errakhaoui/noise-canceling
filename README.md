# noise-canceling

Real-time noise cancellation application for Go using RNNoise (Recurrent Neural Network Noise Suppression).

## Features

- Real-time microphone input capture and processing
- Live noise cancellation using RNNoise ML model
- Real-time audio playback with processed audio
- **Live toggle**: Turn noise cancellation ON/OFF while running
- Thread-safe toggle operations
- High-performance audio processing (48kHz sample rate)

## Requirements

- Go 1.22.3 or higher
- **rnnoise** library (C library)
- **PortAudio** library (for microphone input)
- Audio input/output devices

### Installing Dependencies

#### macOS
```bash
brew install rnnoise portaudio
```

#### Ubuntu/Debian
```bash
sudo apt-get install librnnoise-dev portaudio19-dev
```

#### Arch Linux
```bash
sudo pacman -S rnnoise portaudio
```

## Installation

```bash
# Clone the repository
git clone https://github.com/errakhaoui/noise-canceling.git
cd noise-canceling

# Download Go dependencies
go mod download
```

## Build Commands

### Build the executable
```bash
go build -o noise-canceling example.go
```

### Build with optimizations
```bash
go build -ldflags="-s -w" -o noise-canceling example.go
```

### Build for specific platform
```bash
# For Linux
GOOS=linux GOARCH=amd64 go build -o noise-canceling example.go

# For macOS
GOOS=darwin GOARCH=arm64 go build -o noise-canceling example.go
```

## Running the Application

### Run directly with Go
```bash
go run example.go
```

### Run the built executable
```bash
./noise-canceling
```

### Usage

Once the application is running:

1. The program will start capturing audio from your default microphone
2. Noise cancellation is **ENABLED** by default
3. Processed audio is played back in real-time
4. To toggle noise cancellation:
   - Type `t` and press **Enter** to toggle ON/OFF
   - You'll see: `[TOGGLED] Noise cancellation: ENABLED` or `DISABLED`
5. Press `Ctrl+C` to stop the application

## Testing Commands

### Run all tests
```bash
go test ./...
```

### Run tests with verbose output
```bash
go test ./... -v
```

### Run tests excluding hardware-dependent tests
```bash
go test ./... -short
```

### Run tests with verbose output (short mode)
```bash
go test ./... -short -v
```

### Run tests with race detector
```bash
go test ./... -race
```

### Run tests with coverage
```bash
go test ./... -cover
```

### Generate coverage report
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Run specific package tests
```bash
# Test noise canceller only
go test ./noise_canceller/... -v

# Test input (microphone) only
go test ./input/... -v

# Test output (player) only
go test ./output/... -v
```

## Benchmark Commands

### Run all benchmarks
```bash
go test ./... -bench=.
```

### Run benchmarks with memory statistics
```bash
go test ./noise_canceller/... -bench=. -benchmem
```

### Run specific benchmarks
```bash
# Benchmark toggle operation
go test ./noise_canceller/... -bench=BenchmarkToggle

# Benchmark execute when enabled
go test ./noise_canceller/... -bench=BenchmarkExecuteEnabled

# Benchmark execute when disabled
go test ./noise_canceller/... -bench=BenchmarkExecuteDisabled
```

### Run benchmarks with custom duration
```bash
go test ./noise_canceller/... -bench=. -benchtime=10s
```

## Code Quality Commands

### Run linter
```bash
golangci-lint run
```

### Format code
```bash
go fmt ./...
```

### Check for potential issues
```bash
go vet ./...
```

## Development Commands

### Download dependencies
```bash
go mod download
```

### Update dependencies
```bash
go get -u ./...
go mod tidy
```

### Verify dependencies
```bash
go mod verify
```

### View module information
```bash
go list -m all
```

## Project Structure

```
noise-canceling/
├── example.go              # Main application entry point
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
├── README.md               # This file
├── input/
│   ├── microphone.go       # Microphone input capture
│   └── microphone_test.go  # Tests for microphone
├── noise_canceller/
│   ├── noise_canceller.go      # RNNoise integration
│   └── noise_canceller_test.go # Tests for noise canceller
└── output/
    └── player.go           # Audio playback
```

## How It Works

1. **Input**: Audio is captured from the microphone using PortAudio at 48kHz
2. **Processing**: Audio frames (480 samples) are processed through RNNoise
   - When **ENABLED**: RNNoise ML model removes background noise
   - When **DISABLED**: Audio passes through unmodified
3. **Output**: Processed audio is played back in real-time using Oto
4. **Toggle**: Press 't' + Enter to switch between enabled/disabled states

## Performance

Benchmark results (Apple M1 Pro):

| Operation | Time | Memory |
|-----------|------|--------|
| Toggle() | ~6.7 ns/op | 0 allocs |
| IsEnabled() | ~0.6 ns/op | 0 allocs |
| Execute (disabled) | ~2.2 ns/op | 0 allocs |
| Execute (enabled) | ~44 μs/op | 2 KB/op |

## Technical Details

- **Sample Rate**: 48000 Hz
- **Frame Size**: 480 samples
- **Audio Format**: 16-bit signed integer (int16)
- **Channels**: Mono (1 channel)
- **Buffer Size**: 1024 samples
- **Processing**: Real-time with minimal latency

## Thread Safety

The toggle functionality uses atomic operations to ensure thread-safe state management across concurrent goroutines.


## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
