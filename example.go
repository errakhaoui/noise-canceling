package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/errakhaoui/noise-canceling/input"
	"github.com/errakhaoui/noise-canceling/noise_canceller"
	"github.com/errakhaoui/noise-canceling/output"
	"github.com/gordonklaus/portaudio"
)

func main() {
	// Command-line flags
	listDevices := flag.Bool("list-devices", false, "List all available output devices and exit")
	deviceName := flag.String("device", "", "Output device name (e.g., 'BlackHole' or 'BlackHole 2ch')")
	monitorDevice := flag.String("monitor-device", "", "Additional output device for monitoring (e.g., 'Headphones')")
	flag.Parse()

	// If list-devices flag is set, print devices and exit
	if *listDevices {
		output.PrintAvailableDevices()
		return
	}

	log.Println("Start noise cancellation ...")
	log.Println("Press 't' + Enter to toggle noise cancellation ON/OFF")
	log.Printf("Noise cancellation: ENABLED")

	// Initialize microphone input
	input.StartMicAcquisition()

	// Initialize output device(s)
	var devices []*portaudio.DeviceInfo

	if *deviceName != "" {
		// Try to find and use the specified device
		device, err := output.FindDeviceByName(*deviceName)
		if err != nil {
			log.Fatalf("Error finding device '%s': %v\nRun with -list-devices to see available devices", *deviceName, err)
		}
		devices = append(devices, device)
	}

	if *monitorDevice != "" {
		// Try to find and use the monitor device
		device, err := output.FindDeviceByName(*monitorDevice)
		if err != nil {
			log.Fatalf("Error finding monitor device '%s': %v\nRun with -list-devices to see available devices", *monitorDevice, err)
		}
		devices = append(devices, device)
	}

	// Start output streams
	if err := output.StartOutputStreamToDevices(devices); err != nil {
		log.Fatal(err)
	}

	log.Println("Ready! Audio processing started.")

	// Set up signal handler for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("\nShutting down...")
		input.Close()
		output.Close()
		input.Terminate()
		output.Terminate()
		noise_canceller.Terminate()
		os.Exit(0)
	}()

	// Start keyboard listener in a separate goroutine
	go keyboardListener()

	for {
		// Read audio from the input stream
		input.ReadStream()
		noise_canceller.Execute(input.InputBuffer)
		output.ReadStream(input.InputBuffer)
	}
}

func keyboardListener() {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		if input == "t" {
			newState := noise_canceller.Toggle()
			if newState {
				fmt.Println("\n[TOGGLED] Noise cancellation: ENABLED")
			} else {
				fmt.Println("\n[TOGGLED] Noise cancellation: DISABLED")
			}
		}
	}
}
