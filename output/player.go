package output

import (
	"fmt"
	"log"
	"strings"

	"github.com/gordonklaus/portaudio"
)

const (
	sampleRate   = 48000
	frameSize    = 480
	channelCount = 1
)

var outputStreams []*portaudio.Stream
var outputBuffers [][]int16

func init() {
	if err := portaudio.Initialize(); err != nil {
		log.Fatal(err)
	}
}

// ListOutputDevices lists all available output devices
func ListOutputDevices() ([]*portaudio.DeviceInfo, error) {
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

// PrintAvailableDevices prints all available output devices
func PrintAvailableDevices() {
	devices, err := ListOutputDevices()
	if err != nil {
		log.Printf("Error listing devices: %v", err)
		return
	}

	fmt.Println("\n=== Available Output Devices ===")
	for i, device := range devices {
		fmt.Printf("[%d] %s (Channels: %d)\n", i, device.Name, device.MaxOutputChannels)
	}
	fmt.Println("================================")
}

// FindDeviceByName finds an output device by name (case-insensitive, partial match)
func FindDeviceByName(name string) (*portaudio.DeviceInfo, error) {
	devices, err := ListOutputDevices()
	if err != nil {
		return nil, err
	}

	searchName := strings.ToLower(name)
	for _, device := range devices {
		if strings.Contains(strings.ToLower(device.Name), searchName) {
			return device, nil
		}
	}
	return nil, fmt.Errorf("device containing '%s' not found", name)
}

// StartOutputStream initializes the output stream with default device
func StartOutputStream() error {
	return StartOutputStreamToDevice(nil)
}

// StartOutputStreamToDevice initializes the output stream with a specific device
func StartOutputStreamToDevice(device *portaudio.DeviceInfo) error {
	devices := []*portaudio.DeviceInfo{}
	if device != nil {
		devices = append(devices, device)
	}
	return StartOutputStreamToDevices(devices)
}

// StartOutputStreamToDevices initializes output streams to multiple devices
func StartOutputStreamToDevices(devices []*portaudio.DeviceInfo) error {
	// If no devices specified, use default
	if len(devices) == 0 {
		defaultDevice, err := portaudio.DefaultOutputDevice()
		if err != nil {
			return fmt.Errorf("error getting default output device: %v", err)
		}
		devices = append(devices, defaultDevice)
	}

	// Create a stream for each device
	for _, device := range devices {
		buffer := make([]int16, frameSize)

		var streamParams portaudio.StreamParameters
		streamParams.Output.Channels = channelCount
		streamParams.SampleRate = sampleRate
		streamParams.FramesPerBuffer = frameSize
		streamParams.Output.Device = device

		// Use higher latency for virtual audio devices to prevent underflow
		// BlackHole and similar virtual devices benefit from higher latency
		if device.DefaultHighOutputLatency > 0 {
			streamParams.Output.Latency = device.DefaultHighOutputLatency
		} else {
			streamParams.Output.Latency = device.DefaultLowOutputLatency
		}

		stream, err := portaudio.OpenStream(streamParams, buffer)
		if err != nil {
			// Clean up already opened streams
			closeAllStreams()
			return fmt.Errorf("error opening output stream to %s: %v", device.Name, err)
		}

		err = stream.Start()
		if err != nil {
			_ = stream.Close() // Ignore error on cleanup
			closeAllStreams()
			return fmt.Errorf("error starting output stream to %s: %v", device.Name, err)
		}

		outputStreams = append(outputStreams, stream)
		outputBuffers = append(outputBuffers, buffer)
		log.Printf("Opened output stream to device: %s (latency: %.2fms)",
			device.Name, streamParams.Output.Latency.Seconds()*1000)
	}

	return nil
}

// closeAllStreams is a helper to close all streams (used internally)
func closeAllStreams() {
	for _, stream := range outputStreams {
		if stream != nil {
			_ = stream.Stop()  // Ignore error on cleanup
			_ = stream.Close() // Ignore error on cleanup
		}
	}
	outputStreams = nil
	outputBuffers = nil
}

// ReadStream writes audio data to all output streams
func ReadStream(audioStream []int16) {
	// Validate input size
	if len(audioStream) != frameSize {
		log.Printf("Warning: audio stream size mismatch: expected %d, got %d", frameSize, len(audioStream))
		return
	}

	// Write to all output streams
	for i, stream := range outputStreams {
		// Ensure buffer is the correct size
		if len(outputBuffers[i]) != frameSize {
			log.Printf("Warning: output buffer %d size mismatch", i)
			continue
		}

		// Copy audio data to this stream's buffer
		copy(outputBuffers[i], audioStream)

		// Write to this output stream
		err := stream.Write()
		if err != nil {
			// Underflow errors are common with virtual audio devices and can be ignored
			// They happen when the output can't keep up with the input rate
			errStr := err.Error()
			if errStr != "Output underflowed" {
				log.Printf("Error writing to output stream %d: %v", i, err)
			}
			// Otherwise silently ignore underflow errors
		}
	}
}

// Close closes all output streams
func Close() {
	for i, stream := range outputStreams {
		if stream != nil {
			if err := stream.Stop(); err != nil {
				log.Printf("Error stopping output stream %d: %v", i, err)
			}
			if err := stream.Close(); err != nil {
				log.Printf("Error closing output stream %d: %v", i, err)
			}
		}
	}
	outputStreams = nil
	outputBuffers = nil
}

// Terminate closes PortAudio completely (call only on final exit)
func Terminate() {
	if err := portaudio.Terminate(); err != nil {
		log.Printf("Error terminating portaudio: %v", err)
	}
}
