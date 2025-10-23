package gui

import (
	"fmt"
	"image/color"
	"log"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/errakhaoui/noise-canceling/input"
	"github.com/errakhaoui/noise-canceling/noise_canceller"
	"github.com/errakhaoui/noise-canceling/output"
	"github.com/gordonklaus/portaudio"
)

// AudioProcessor manages the audio processing state
type AudioProcessor struct {
	running            bool
	stopChan           chan bool
	mu                 sync.Mutex
	inputDeviceIndex   int
	outputDeviceIndex  int
	monitorDeviceIndex int
	noiseCancelEnabled bool
}

var processor = &AudioProcessor{
	running:            false,
	stopChan:           make(chan bool),
	noiseCancelEnabled: true,
	inputDeviceIndex:   -1,
	outputDeviceIndex:  -1,
	monitorDeviceIndex: -1,
}

// getInputDevices returns all available input devices
func getInputDevices() ([]*portaudio.DeviceInfo, error) {
	if err := portaudio.Initialize(); err != nil {
		return nil, err
	}
	defer func() {
		_ = portaudio.Terminate() // Ignore error as this is temporary initialization
	}()

	devices, err := portaudio.Devices()
	if err != nil {
		return nil, err
	}

	var inputDevices []*portaudio.DeviceInfo
	for _, device := range devices {
		if device.MaxInputChannels > 0 {
			inputDevices = append(inputDevices, device)
		}
	}
	return inputDevices, nil
}

// getOutputDevices returns all available output devices
func getOutputDevices() ([]*portaudio.DeviceInfo, error) {
	if err := portaudio.Initialize(); err != nil {
		return nil, err
	}
	defer func() {
		_ = portaudio.Terminate() // Ignore error as this is temporary initialization
	}()

	devices, err := portaudio.Devices()
	if err != nil {
		return nil, err
	}

	var outputDevices []*portaudio.DeviceInfo
	for _, device := range devices {
		if device.MaxOutputChannels > 0 {
			outputDevices = append(outputDevices, device)
		}
	}
	return outputDevices, nil
}

// startAudioProcessing starts the audio processing loop
func startAudioProcessing(inputDevices, outputDevices []*portaudio.DeviceInfo) error {
	processor.mu.Lock()
	if processor.running {
		processor.mu.Unlock()
		return fmt.Errorf("already running")
	}
	processor.running = true
	processor.mu.Unlock()

	// Initialize input
	input.StartMicAcquisition()

	// Initialize output devices
	var devicesToUse []*portaudio.DeviceInfo
	if processor.outputDeviceIndex >= 0 && processor.outputDeviceIndex < len(outputDevices) {
		devicesToUse = append(devicesToUse, outputDevices[processor.outputDeviceIndex])
	}
	if processor.monitorDeviceIndex >= 0 && processor.monitorDeviceIndex < len(outputDevices) {
		devicesToUse = append(devicesToUse, outputDevices[processor.monitorDeviceIndex])
	}

	if err := output.StartOutputStreamToDevices(devicesToUse); err != nil {
		input.Close()
		processor.mu.Lock()
		processor.running = false
		processor.mu.Unlock()
		return err
	}

	// Set initial noise cancellation state
	if processor.noiseCancelEnabled {
		noise_canceller.Enable()
	} else {
		noise_canceller.Disable()
	}

	// Start processing loop in a goroutine
	go func() {
		for {
			select {
			case <-processor.stopChan:
				return
			default:
				// Read audio from the input stream
				input.ReadStream()
				noise_canceller.Execute(input.InputBuffer)
				output.ReadStream(input.InputBuffer)
			}
		}
	}()

	return nil
}

// stopAudioProcessing stops the audio processing
func stopAudioProcessing() {
	processor.mu.Lock()
	defer processor.mu.Unlock()

	if !processor.running {
		return
	}

	processor.stopChan <- true
	processor.running = false

	// Close audio streams (but keep RNNoise state alive for restart)
	input.Close()
	output.Close()
	// Note: noise_canceller is NOT closed here - it persists across start/stop cycles
}

// toggleNoiseCancellation toggles noise cancellation on/off
func toggleNoiseCancellation(enabled bool) {
	processor.mu.Lock()
	processor.noiseCancelEnabled = enabled
	processor.mu.Unlock()

	if enabled {
		noise_canceller.Enable()
	} else {
		noise_canceller.Disable()
	}
}

// CreateGUI creates and displays the main GUI window
func CreateGUI() {
	myApp := app.New()
	myWindow := myApp.NewWindow("ClearVox")
	myWindow.Resize(fyne.NewSize(450, 400))

	// Get available devices
	inputDevices, err := getInputDevices()
	if err != nil {
		log.Fatalf("Failed to get input devices: %v", err)
	}

	outputDevices, err := getOutputDevices()
	if err != nil {
		log.Fatalf("Failed to get output devices: %v", err)
	}

	// Create device name lists for dropdowns
	inputDeviceNames := []string{}
	for _, dev := range inputDevices {
		inputDeviceNames = append(inputDeviceNames, dev.Name)
	}

	outputDeviceNames := []string{}
	for _, dev := range outputDevices {
		outputDeviceNames = append(outputDeviceNames, dev.Name)
	}

	monitorDeviceNames := []string{"None"}
	monitorDeviceNames = append(monitorDeviceNames, outputDeviceNames...)

	// Set default selections
	defaultInputIdx := 0
	for i, dev := range inputDevices {
		if dev.MaxInputChannels > 0 {
			defaultDevice, _ := portaudio.DefaultInputDevice()
			if defaultDevice != nil && dev.Name == defaultDevice.Name {
				defaultInputIdx = i
				break
			}
		}
	}

	defaultOutputIdx := 0
	for i, dev := range outputDevices {
		if dev.MaxOutputChannels > 0 {
			defaultDevice, _ := portaudio.DefaultOutputDevice()
			if defaultDevice != nil && dev.Name == defaultDevice.Name {
				defaultOutputIdx = i
				break
			}
		}
	}

	processor.inputDeviceIndex = defaultInputIdx
	processor.outputDeviceIndex = defaultOutputIdx
	processor.monitorDeviceIndex = -1 // None

	// Create widgets
	inputLabel := widget.NewLabel("Input Device:")
	inputSelect := widget.NewSelect(inputDeviceNames, func(value string) {
		for i, name := range inputDeviceNames {
			if name == value {
				processor.inputDeviceIndex = i
				break
			}
		}
	})
	if len(inputDeviceNames) > 0 {
		inputSelect.SetSelected(inputDeviceNames[defaultInputIdx])
	}

	outputLabel := widget.NewLabel("Output Device (Virtual Mic):")
	outputSelect := widget.NewSelect(outputDeviceNames, func(value string) {
		for i, name := range outputDeviceNames {
			if name == value {
				processor.outputDeviceIndex = i
				break
			}
		}
	})
	if len(outputDeviceNames) > 0 {
		outputSelect.SetSelected(outputDeviceNames[defaultOutputIdx])
	}

	monitorLabel := widget.NewLabel("Monitor Device (optional):")
	monitorSelect := widget.NewSelect(monitorDeviceNames, func(value string) {
		if value == "None" {
			processor.monitorDeviceIndex = -1
		} else {
			for i, name := range outputDeviceNames {
				if name == value {
					processor.monitorDeviceIndex = i
					break
				}
			}
		}
	})
	monitorSelect.SetSelected("None")

	noiseCancelCheck := widget.NewCheck("Enable Noise Cancellation", func(checked bool) {
		toggleNoiseCancellation(checked)
	})
	noiseCancelCheck.SetChecked(true)

	// Status indicator
	statusCircle := canvas.NewCircle(color.NRGBA{R: 255, G: 0, B: 0, A: 255}) // Red = stopped
	statusCircle.Resize(fyne.NewSize(15, 15))
	statusLabel := widget.NewLabel("Status: Stopped")
	statusContainer := container.NewHBox(statusCircle, statusLabel)

	// Start/Stop buttons
	startButton := widget.NewButton("Start", nil)
	stopButton := widget.NewButton("Stop", nil)
	stopButton.Disable()

	startButton.OnTapped = func() {
		err := startAudioProcessing(inputDevices, outputDevices)
		if err != nil {
			log.Printf("Error starting audio processing: %v", err)
			statusLabel.SetText(fmt.Sprintf("Status: Error - %v", err))
			return
		}

		statusCircle.FillColor = color.NRGBA{R: 0, G: 255, B: 0, A: 255} // Green = running
		statusCircle.Refresh()
		statusLabel.SetText("Status: Running")
		startButton.Disable()
		stopButton.Enable()
		inputSelect.Disable()
		outputSelect.Disable()
		monitorSelect.Disable()
	}

	stopButton.OnTapped = func() {
		stopAudioProcessing()

		statusCircle.FillColor = color.NRGBA{R: 255, G: 0, B: 0, A: 255} // Red = stopped
		statusCircle.Refresh()
		statusLabel.SetText("Status: Stopped")
		startButton.Enable()
		stopButton.Disable()
		inputSelect.Enable()
		outputSelect.Enable()
		monitorSelect.Enable()
	}

	buttonContainer := container.NewHBox(startButton, stopButton)

	// Layout
	content := container.NewVBox(
		widget.NewLabelWithStyle("ClearVox", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		inputLabel,
		inputSelect,
		outputLabel,
		outputSelect,
		monitorLabel,
		monitorSelect,
		widget.NewSeparator(),
		noiseCancelCheck,
		widget.NewSeparator(),
		buttonContainer,
		widget.NewSeparator(),
		statusContainer,
	)

	myWindow.SetContent(content)

	// Set cleanup function when window closes
	myWindow.SetOnClosed(func() {
		// Stop audio processing if running
		stopAudioProcessing()

		// Terminate all resources completely
		input.Terminate()
		output.Terminate()
		noise_canceller.Terminate()

		// Cleanup any temp files
		cleanupTempFiles()
	})

	// Show window first
	myWindow.Show()

	// Check if BlackHole is installed and offer to install it if not
	go func() {
		installed, err := isBlackHoleInstalled()
		if err != nil {
			log.Printf("Error checking BlackHole installation: %v", err)
			return
		}

		if !installed {
			// Show setup dialog
			dialog.ShowConfirm(
				"ClearVox Virtual Microphone Setup",
				"To use ClearVox as a virtual microphone in video calls, you need the BlackHole audio driver.\n\n"+
					"Would you like to install it now?\n\n"+
					"This will:\n"+
					"• Download BlackHole 2ch (~1 MB)\n"+
					"• Install the audio driver (requires admin password)\n"+
					"• Enable virtual microphone for video calls",
				func(install bool) {
					if install {
						handleBlackHoleInstallation(myWindow)
					}
				},
				myWindow,
			)
		}
	}()

	myApp.Run()
}

// handleBlackHoleInstallation downloads and installs BlackHole
func handleBlackHoleInstallation(window fyne.Window) {
	// Show progress dialog
	progressBar := widget.NewProgressBar()
	progressLabel := widget.NewLabel("Downloading BlackHole...")

	progressContent := container.NewVBox(
		progressLabel,
		progressBar,
	)

	progressDialog := dialog.NewCustom("Installing Virtual Microphone", "Cancel", progressContent, window)
	progressDialog.Show()

	go func() {
		// Download BlackHole
		pkgPath, err := downloadBlackHole(func(downloaded, total int64) {
			if total > 0 {
				progress := float64(downloaded) / float64(total)
				progressBar.SetValue(progress)
				progressLabel.SetText(fmt.Sprintf("Downloading... %.1f%%", progress*100))
			}
		})

		if err != nil {
			progressDialog.Hide()
			dialog.ShowError(
				fmt.Errorf("failed to download BlackHole: %w", err),
				window,
			)
			return
		}

		progressLabel.SetText("Opening installer...")
		progressBar.SetValue(1.0)

		// Install BlackHole
		if err := installBlackHole(pkgPath); err != nil {
			progressDialog.Hide()
			dialog.ShowError(
				fmt.Errorf("failed to open installer: %w", err),
				window,
			)
			return
		}

		progressDialog.Hide()

		// Show success message
		dialog.ShowInformation(
			"Installation Started",
			"The BlackHole installer has been opened.\n\n"+
				"Please:\n"+
				"1. Follow the installer prompts\n"+
				"2. Enter your admin password when asked\n"+
				"3. Complete the installation\n\n"+
				"Once installed, restart ClearVox to use the Virtual Microphone feature.",
			window,
		)
	}()
}
