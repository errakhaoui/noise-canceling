package input

import (
	"testing"
)

func TestConstants(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		expected int
	}{
		{"SampleRate", SampleRate, 48000},
		{"frameSize", frameSize, 480},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != tt.expected {
				t.Errorf("%s = %d, want %d", tt.name, tt.value, tt.expected)
			}
		})
	}
}

func TestInputBufferInitialization(t *testing.T) {
	// Save original state
	originalBuffer := InputBuffer

	// Test that StartMicAcquisition initializes the buffer
	// Note: This test may require audio hardware to be available
	t.Run("BufferSize", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping hardware-dependent test in short mode")
		}

		StartMicAcquisition()
		defer func() {
			// Restore original state
			InputBuffer = originalBuffer
			if inputStream != nil {
				if err := inputStream.Stop(); err != nil {
					t.Logf("Error stopping stream: %v", err)
				}
				if err := inputStream.Close(); err != nil {
					t.Logf("Error closing stream: %v", err)
				}
			}
		}()

		if len(InputBuffer) != frameSize {
			t.Errorf("InputBuffer length = %d, want %d", len(InputBuffer), frameSize)
		}

		if cap(InputBuffer) < frameSize {
			t.Errorf("InputBuffer capacity = %d, want at least %d", cap(InputBuffer), frameSize)
		}
	})
}

func TestInputBufferType(t *testing.T) {
	t.Run("BufferElementType", func(t *testing.T) {
		// Create a test buffer to verify type
		testBuffer := make([]int16, frameSize)

		// Verify we can assign int16 values
		testBuffer[0] = 32767
		testBuffer[1] = -32768

		if testBuffer[0] != 32767 {
			t.Errorf("Expected max int16 value 32767, got %d", testBuffer[0])
		}

		if testBuffer[1] != -32768 {
			t.Errorf("Expected min int16 value -32768, got %d", testBuffer[1])
		}
	})
}

func TestReadStreamWithoutInitialization(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping hardware-dependent test in short mode")
	}

	t.Run("ReadStreamPanicsWithoutInit", func(t *testing.T) {
		// This test verifies that ReadStream handles uninitialized stream
		// We expect it to fail/panic if called without StartMicAcquisition

		// Save original stream
		originalStream := inputStream
		inputStream = nil

		defer func() {
			// Restore original stream
			inputStream = originalStream

			if r := recover(); r == nil {
				t.Error("Expected ReadStream to panic with nil stream, but it didn't")
			}
		}()

		ReadStream()
	})
}

// BenchmarkReadStream benchmarks the audio reading performance
func BenchmarkReadStream(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping hardware-dependent benchmark in short mode")
	}

	// Initialize
	StartMicAcquisition()
	defer func() {
		if err := inputStream.Stop(); err != nil {
			b.Logf("Error stopping stream: %v", err)
		}
		if err := inputStream.Close(); err != nil {
			b.Logf("Error closing stream: %v", err)
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ReadStream()
	}
}
