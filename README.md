# noise-canceling

Real-time noise cancellation application for Go using RNNoise (Recurrent Neural Network Noise Suppression).

Available in two versions:
- **GUI Application**: User-friendly graphical interface (recommended)
- **CLI Application**: Command-line interface with keyboard controls

## Features

- Real-time microphone input capture and processing
- Live noise cancellation using RNNoise ML model
- Real-time audio playback with processed audio
- **Virtual microphone support**: Output to virtual audio devices (BlackHole) for use in Google Meet, Zoom, etc.
- **GUI Application**: Easy-to-use graphical interface with device selection and controls
- **Live toggle**: Turn noise cancellation ON/OFF while running
- Thread-safe toggle operations
- High-performance audio processing (48kHz sample rate)
- Device selection: Choose input, output, and monitor devices
- Multi-device output: Send audio to multiple devices simultaneously

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

**For Virtual Microphone Support (Optional):**
Install BlackHole to create a virtual audio device that can be used as a microphone in apps like Google Meet, Zoom, Discord, etc.

```bash
brew install blackhole-2ch
```

After installation, BlackHole will appear as an audio device named "BlackHole 2ch" in your system's audio settings and in applications.

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

### Build the GUI application (recommended)
```bash
go build -o noise-canceling-gui gui_main.go
```

### Build the CLI application
```bash
go build -o noise-canceling example.go
```

### Build both versions with optimizations
```bash
# GUI version
go build -ldflags="-s -w" -o noise-canceling-gui gui_main.go

# CLI version
go build -ldflags="-s -w" -o noise-canceling example.go
```

### Build for specific platform
```bash
# GUI - For macOS
GOOS=darwin GOARCH=arm64 go build -o noise-canceling-gui gui_main.go

# CLI - For macOS
GOOS=darwin GOARCH=arm64 go build -o noise-canceling example.go
```

## Running the Application

### GUI Application (Recommended)

#### Run directly with Go
```bash
go run gui_main.go
```

#### Run the built executable
```bash
./noise-canceling-gui
```

#### GUI Interface

The graphical interface provides an intuitive way to configure and control noise cancellation:

```
┌─────────────────────────────────────────┐
│  Noise Cancellation Control             │
├─────────────────────────────────────────┤
│                                         │
│  Input Device:                          │
│  ┌───────────────────────────────────┐ │
│  │ Built-in Microphone              ▼│ │
│  └───────────────────────────────────┘ │
│                                         │
│  Output Device:                         │
│  ┌───────────────────────────────────┐ │
│  │ BlackHole 2ch                    ▼│ │
│  └───────────────────────────────────┘ │
│                                         │
│  Monitor Device (optional):             │
│  ┌───────────────────────────────────┐ │
│  │ Écouteurs externes               ▼│ │
│  └───────────────────────────────────┘ │
│                                         │
│  ☑ Enable Noise Cancellation           │
│                                         │
│  [ Start ]  [ Stop ]                    │
│                                         │
│  Status: ● Running                      │
└─────────────────────────────────────────┘
```

> **Note:** To add an actual screenshot here, run the GUI and take a screenshot:
> ```bash
> # On macOS, press Cmd+Shift+4, then Space, then click the GUI window
> # Save as: docs/gui-screenshot.png
> ```
> Then add to README: `![GUI Screenshot](docs/gui-screenshot.png)`

The GUI provides:
- **Input Device**: Select your physical microphone
- **Output Device**: Select where to send processed audio (e.g., BlackHole for virtual mic)
- **Monitor Device**: Optionally select your headphones/speakers to hear the output
- **Noise Cancellation**: Toggle on/off with checkbox (works in real-time)
- **Start/Stop**: Control audio processing with buttons
- **Status Indicator**: Visual feedback (red = stopped, green = running)

---

## Virtual Microphone Setup for Google Meet/Zoom

### How It Works - Audio Flow Diagram

```
┌──────────────────────────────────────────────────────────────┐
│                    Your Computer                             │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  [Physical Microphone] ──────┐                              │
│   (Built-in or External)     │                              │
│                              │                              │
│                              ▼                              │
│                      ┌──────────────┐                       │
│                      │   GUI App    │                       │
│                      │ ┌──────────┐ │                       │
│                      │ │ RNNoise  │ │ ← Removes noise      │
│                      │ │Processor │ │   in real-time       │
│                      │ └──────────┘ │                       │
│                      └──────────────┘                       │
│                              │                              │
│                        ┌─────┴─────┐                        │
│                        │           │                        │
│                        ▼           ▼                        │
│                [BlackHole 2ch] [Headphones]                 │
│                 (Virtual Mic)  (Monitor)                    │
│                        │                                    │
│                        │                                    │
└────────────────────────┼────────────────────────────────────┘
                         │
                         ▼
                  ┌─────────────┐
                  │ Google Meet │ ← Selects BlackHole
                  │  Zoom, etc  │   as microphone
                  └─────────────┘
```

### Step-by-Step Configuration

#### Step 1: Install BlackHole (One-time setup)
```bash
brew install blackhole-2ch
```

#### Step 2: Launch the GUI Application
```bash
./noise-canceling-gui
```

#### Step 3: Configure Devices in the GUI

**Input Device:**
- Select: **"Built-in Microphone"** (or your physical microphone name)
- This is where your voice comes FROM

**Output Device:**
- Select: **"BlackHole 2ch"**
- This is your VIRTUAL microphone (where processed audio goes TO)
- BlackHole acts as a bridge between your app and Google Meet

**Monitor Device (optional but recommended):**
- Select: **"Écouteurs externes"** or **"Haut-parleurs MacBook Pro"**
- This lets YOU hear what you're sending to others
- Without this, you won't hear yourself speaking

**Enable Noise Cancellation:**
- Check: **☑ Enable Noise Cancellation**
- You can toggle this on/off even while running

#### Step 4: Click Start
- Status indicator will turn GREEN
- Audio processing begins

#### Step 5: Configure Google Meet/Zoom

**In Google Meet:**
1. Click the three dots (⋮) → Settings
2. Go to **Audio** tab
3. **Microphone**: Select **"BlackHole 2ch"**
4. **Speakers**: Keep your normal speakers/headphones (NOT BlackHole)

**In Zoom:**
1. Click Settings (⚙️) → Audio
2. **Microphone**: Select **"BlackHole 2ch"**
3. **Speaker**: Keep your normal speakers/headphones (NOT BlackHole)

### Why This Works

**BlackHole is Both an Input AND Output Device:**
- In YOUR app: BlackHole = **Output** (where you send processed audio)
- In Google Meet: BlackHole = **Input/Microphone** (where Meet reads from)
- BlackHole passes audio between the two applications

**The Complete Audio Path:**
1. Your physical mic captures your voice
2. GUI app receives raw audio
3. RNNoise removes background noise
4. Clean audio is sent to BlackHole (output)
5. Google Meet reads clean audio from BlackHole (input)
6. You hear yourself via Monitor Device (headphones)
7. Other people hear your noise-canceled voice

### Quick Configuration Summary

**Your GUI App Configuration:**
```
Input Device:    Built-in Microphone
Output Device:   BlackHole 2ch
Monitor Device:  Écouteurs externes (optional)
☑ Enable Noise Cancellation
[Start]
```

**Google Meet/Zoom Configuration:**
```
Microphone: BlackHole 2ch
Speakers:   Your normal speakers (NOT BlackHole)
```

**Result:** People in your meeting hear your noise-canceled audio! 🎉

### CLI Application

#### Run directly with Go
```bash
go run example.go
```

#### Run the built executable
```bash
./noise-canceling
```

#### List Available Audio Devices
```bash
./noise-canceling -list-devices
```

#### Run with Virtual Microphone (BlackHole)
```bash
./noise-canceling -device blackhole
```

#### Run with Virtual Microphone + Headphone Monitoring
```bash
./noise-canceling -device blackhole -monitor-device headphones
```
Or for external headphones:
```bash
./noise-canceling -device blackhole -monitor-device "écouteurs"
```

### CLI Usage Details

#### Basic Usage (Speakers/Headphones)

Run without any flags to output processed audio to your default speakers:

```bash
./noise-canceling
```

1. The program will start capturing audio from your default microphone
2. Noise cancellation is **ENABLED** by default
3. Processed audio is played back in real-time through speakers
4. To toggle noise cancellation:
   - Type `t` and press **Enter** to toggle ON/OFF
   - You'll see: `[TOGGLED] Noise cancellation: ENABLED` or `DISABLED`
5. Press `Ctrl+C` to stop the application

#### Virtual Microphone Usage (For Google Meet, Zoom, etc.)

To use noise-canceled audio as a virtual microphone in video conferencing apps:

**Step 1:** Install BlackHole (see installation instructions above)

**Step 2:** List available devices to verify BlackHole is installed:
```bash
./noise-canceling -list-devices
```
You should see output like:
```
=== Available Output Devices ===
[0] BlackHole 2ch (Channels: 2)
[1] Écouteurs externes (Channels: 2)
[2] Haut-parleurs MacBook Pro (Channels: 2)
================================
```

**Step 3:** Run the application with BlackHole AND your headphones for monitoring:
```bash
./noise-canceling -device blackhole -monitor-device "écouteurs"
```
Or adjust based on your device names:
```bash
# For built-in speakers
./noise-canceling -device blackhole -monitor-device "haut-parleurs"

# The device name matching is case-insensitive and partial
./noise-canceling -device blackhole -monitor-device headphones
```

**Step 4:** In your video conferencing app (Google Meet, Zoom, Discord, etc.):
- Open audio settings
- Select **"BlackHole 2ch"** as your **microphone/input device**
- Your physical microphone audio will now be processed through noise cancellation before being sent to the app
- You'll hear the processed audio in your headphones/speakers

**Step 5:** Toggle noise cancellation on/off anytime by typing `t` + Enter in the terminal

**Note:** If you only use `-device blackhole` without `-monitor-device`, you won't hear the audio yourself. The `-monitor-device` flag allows you to hear what you're sending to the video conferencing app.

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
├── gui_main.go             # GUI application entry point
├── example.go              # CLI application entry point
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
├── README.md               # This file
├── gui/
│   └── gui.go              # GUI components and logic (Fyne)
├── input/
│   ├── microphone.go       # Microphone input capture
│   └── microphone_test.go  # Tests for microphone
├── noise_canceller/
│   ├── noise_canceller.go      # RNNoise integration
│   └── noise_canceller_test.go # Tests for noise canceller
└── output/
    └── player.go           # Audio playback and device management
```

## How It Works

1. **Input**: Audio is captured from the microphone using PortAudio at 48kHz
2. **Processing**: Audio frames (480 samples) are processed through RNNoise
   - When **ENABLED**: RNNoise ML model removes background noise
   - When **DISABLED**: Audio passes through unmodified
3. **Output**: Processed audio is sent to the selected output device using PortAudio
   - **Default mode**: Plays to speakers/headphones for monitoring
   - **Virtual mic mode**: Outputs to BlackHole, making it available as a microphone in other apps
4. **Toggle**: Press 't' + Enter to switch between enabled/disabled states

## Command-Line Options

```bash
./noise-canceling [OPTIONS]
```

**Options:**
- `-list-devices` : List all available output audio devices and exit
- `-device <name>` : Specify primary output device name (e.g., "blackhole", "BlackHole 2ch")
- `-monitor-device <name>` : Specify additional output device for monitoring (e.g., "headphones", "speakers")

**Examples:**
```bash
# List all available devices
./noise-canceling -list-devices

# Use default output device (speakers)
./noise-canceling

# Output to BlackHole for virtual microphone (no monitoring)
./noise-canceling -device blackhole

# Output to BlackHole + hear audio in headphones
./noise-canceling -device blackhole -monitor-device headphones

# Output to BlackHole + hear audio in speakers
./noise-canceling -device blackhole -monitor-device speakers

# Output to any specific device
./noise-canceling -device "External Headphones"

# Output to multiple devices simultaneously
./noise-canceling -device "Device 1" -monitor-device "Device 2"
```

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
- **Frame Size**: 480 samples (10ms at 48kHz)
- **Audio Format**: 16-bit signed integer (int16)
- **Channels**: Mono (1 channel)
- **Input Latency**: Low (hardware default)
- **Output Latency**: High for virtual devices, low for physical devices
- **Processing**: Real-time with minimal latency (~10-20ms typical)

### Buffer Management

The application uses careful buffer management to ensure smooth audio processing:

- **Input Buffer**: 480 samples per frame
- **Output Buffer**: 480 samples per frame (matches input)
- **Virtual Devices**: Uses higher latency settings to prevent underflow
- **Physical Devices**: Uses lower latency for minimal delay
- **Error Handling**: Automatically suppresses non-critical underflow errors

## Thread Safety

The toggle functionality uses atomic operations to ensure thread-safe state management across concurrent goroutines.

## Troubleshooting

### Common Issues and Solutions

#### 1. "BlackHole 2ch not showing up in Google Meet"

**Problem:** Google Meet doesn't list BlackHole as a microphone option.

**Solutions:**
- Verify BlackHole is installed: `brew list blackhole-2ch`
- Restart Google Meet (close tab and reopen)
- Restart your browser completely
- Check System Settings → Sound → Input to confirm BlackHole appears

#### 2. "I can't hear myself when using BlackHole"

**Problem:** No audio monitoring when outputting to BlackHole.

**Solution:**
- In the GUI, select your headphones/speakers as **Monitor Device**
- This sends audio to BOTH BlackHole (for Google Meet) AND your headphones (for you to hear)

#### 3. "People in Google Meet can't hear me"

**Checklist:**
- ✅ GUI app shows "Status: Running" (green indicator)
- ✅ Input Device is set to your physical microphone
- ✅ Output Device is set to BlackHole 2ch
- ✅ Google Meet microphone is set to BlackHole 2ch (NOT your physical mic)
- ✅ Check Google Meet hasn't muted you
- ✅ Try stopping and starting the GUI app

#### 4. "Audio sounds robotic or distorted"

**Solutions:**
- Ensure your microphone is selected as Input (not BlackHole)
- Check that Google Meet speaker is NOT set to BlackHole (should be regular speakers)
- Verify BlackHole is only used as output in GUI and input in Google Meet
- Restart the GUI application

#### 5. "Error: Failed to get input/output devices"

**Problem:** PortAudio can't enumerate devices.

**Solutions:**
- Restart the application
- Check system permissions for microphone access (System Settings → Privacy → Microphone)
- Reinstall PortAudio: `brew reinstall portaudio`

#### 6. "GUI app crashes on Start"

**Solutions:**
- Check that devices are properly selected (not empty)
- Ensure another app isn't exclusively using the microphone
- Close other audio apps (Zoom, Skype, etc.) and try again
- Check Console.app for error logs

#### 7. "Noise cancellation doesn't seem to work"

**Checklist:**
- ✅ "Enable Noise Cancellation" checkbox is checked
- ✅ Try toggling it off and on while running
- ✅ RNNoise works best with constant background noise (fans, AC, keyboard typing)
- ✅ It won't remove sudden loud sounds or voices in the background as effectively

#### 8. "Build fails with 'library not found' error"

**Solutions:**
```bash
# Reinstall dependencies
brew reinstall rnnoise portaudio

# Clean and rebuild
go clean
go mod tidy
go build -o noise-canceling-gui gui_main.go
```

#### 9. "Output underflowed" errors in the logs

**Problem:** You see "Output underflowed" errors when writing to BlackHole or other virtual devices.

**Solution:**
- This is NORMAL and expected with virtual audio devices like BlackHole
- These errors are automatically suppressed in the latest version
- They indicate the output buffer temporarily ran out of data, but audio continues normally
- If audio is working fine, you can safely ignore these messages

**Why it happens:**
- Virtual audio devices have different timing characteristics than physical devices
- The app now uses higher latency for virtual devices to minimize these errors
- Occasional underflows don't affect audio quality significantly

### Verifying Your Setup

Run this checklist to verify everything is configured correctly:

**In Your GUI App:**
```
[ ] Input Device = Your physical microphone
[ ] Output Device = BlackHole 2ch
[ ] Monitor Device = Your headphones (optional)
[ ] Enable Noise Cancellation = Checked
[ ] Status = Running (green)
```

**In Google Meet:**
```
[ ] Microphone = BlackHole 2ch
[ ] Speakers = Your normal speakers (NOT BlackHole)
[ ] Microphone is not muted in Google Meet
```

**Test:**
1. Speak into your microphone
2. You should hear yourself in your headphones (if Monitor Device is set)
3. Check Google Meet's microphone meter - it should show activity when you speak

### Still Having Issues?

1. Check the application logs in the terminal where you launched the GUI
2. Test with a different application (Discord, Zoom) to isolate the issue
3. Create an issue on GitHub with:
   - Your OS version
   - Error messages from the terminal
   - Steps to reproduce the issue

## FAQ (Frequently Asked Questions)

### Q: Does this work on Windows or Linux?
**A:** The code is cross-platform, but you'll need a BlackHole alternative:
- **Windows:** Use VB-Audio Virtual Cable or VoiceMeeter
- **Linux:** Use PulseAudio's virtual sinks or PipeWire

### Q: Can I use this with Discord, Slack, or other apps?
**A:** Yes! Any application that lets you select an audio input device will work. Just select BlackHole 2ch as the microphone in that app's settings.

### Q: Does noise cancellation affect audio quality?
**A:** RNNoise is designed to preserve voice quality while removing background noise. You might notice a slight change in your voice tone, but most users find it improves overall audio quality by removing distractions.

### Q: How much CPU does this use?
**A:** Very minimal! RNNoise is highly optimized:
- Processing: ~44 microseconds per frame
- CPU usage: typically <5% on modern hardware
- No GPU required

### Q: Can I use my built-in speakers instead of headphones?
**A:** Yes, but be careful of feedback loops:
- If using BlackHole output, you MUST use the Monitor Device to hear yourself
- Your speakers might pick up their own output and create echo
- Headphones are strongly recommended when using virtual microphone mode

### Q: Why is there a delay/latency?
**A:** Audio latency is typically 20-50ms, which is barely noticeable. If you experience significant delay:
- Check your buffer size settings
- Close other audio applications
- Ensure your audio interface drivers are up to date

### Q: Can I run multiple instances?
**A:** No, only one instance can access the microphone at a time. Running multiple instances will cause device access errors.

### Q: Do I need to keep the terminal/GUI window open?
**A:** Yes, the application must keep running while you're in your call. Closing it will stop the noise cancellation.

### Q: Can I save my device preferences?
**A:** Currently, no. Device selection is reset each time you launch the app. This is a potential future enhancement.

### Q: Is my audio data sent anywhere?
**A:** No! All processing happens locally on your computer. Your audio never leaves your machine. The app doesn't make any network connections.

### Q: What's the difference between GUI and CLI versions?
**A:**
- **GUI:** Easier to use, visual device selection, real-time controls
- **CLI:** Lightweight, can be scripted, command-line flags
- Both use the same audio processing engine

### Q: Why do I need BlackHole? Can't the app create a virtual device?
**A:** Creating virtual audio devices on macOS requires system-level drivers. BlackHole is a well-tested, open-source solution that handles this complexity. Building our own would require:
- Code signing
- System extensions
- Kernel-level programming
- Extensive testing across macOS versions

BlackHole (~16KB) is the minimal, trusted solution.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
