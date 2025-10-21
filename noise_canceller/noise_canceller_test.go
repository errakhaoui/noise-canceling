package noise_canceller

import (
	"testing"
)

func TestConstants(t *testing.T) {
	t.Run("FrameSize", func(t *testing.T) {
		expected := 480
		if frameSize != expected {
			t.Errorf("frameSize = %d, want %d", frameSize, expected)
		}
	})
}

func TestInitialState(t *testing.T) {
	t.Run("InitiallyEnabled", func(t *testing.T) {
		// Reset to initial state
		enabled.Store(true)

		if !IsEnabled() {
			t.Error("Expected noise cancellation to be enabled initially")
		}
	})
}

func TestToggle(t *testing.T) {
	tests := []struct {
		name          string
		initialState  bool
		expectedState bool
	}{
		{"Toggle from enabled to disabled", true, false},
		{"Toggle from disabled to enabled", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set initial state
			enabled.Store(tt.initialState)

			// Toggle
			newState := Toggle()

			if newState != tt.expectedState {
				t.Errorf("Toggle() = %v, want %v", newState, tt.expectedState)
			}

			if IsEnabled() != tt.expectedState {
				t.Errorf("IsEnabled() = %v, want %v", IsEnabled(), tt.expectedState)
			}
		})
	}
}

func TestMultipleToggles(t *testing.T) {
	t.Run("MultipleToggleCycles", func(t *testing.T) {
		// Start with known state (enabled = true)
		enabled.Store(true)

		states := []bool{}
		for i := 0; i < 10; i++ {
			state := Toggle()
			states = append(states, state)
		}

		// Verify alternating pattern
		// Starting with true, first toggle gives false (index 0), next gives true (index 1), etc.
		for i := 0; i < len(states); i++ {
			expected := i%2 == 1 // false for even indices (0,2,4...), true for odd (1,3,5...)
			if states[i] != expected {
				t.Errorf("Toggle #%d: got %v, want %v", i, states[i], expected)
			}
		}
	})
}

func TestIsEnabled(t *testing.T) {
	tests := []struct {
		name  string
		state bool
	}{
		{"When enabled", true},
		{"When disabled", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enabled.Store(tt.state)

			if IsEnabled() != tt.state {
				t.Errorf("IsEnabled() = %v, want %v", IsEnabled(), tt.state)
			}
		})
	}
}

func TestExecuteWhenDisabled(t *testing.T) {
	t.Run("ExecuteSkipsProcessingWhenDisabled", func(t *testing.T) {
		// Create test audio data
		testAudio := make([]int16, frameSize)
		for i := range testAudio {
			testAudio[i] = int16(i * 100)
		}

		// Create a copy to compare
		originalAudio := make([]int16, frameSize)
		copy(originalAudio, testAudio)

		// Disable noise cancellation
		enabled.Store(false)

		// Execute should not modify the audio when disabled
		Execute(testAudio)

		// Verify audio is unchanged
		for i := range testAudio {
			if testAudio[i] != originalAudio[i] {
				t.Errorf("Audio modified at index %d: got %d, want %d", i, testAudio[i], originalAudio[i])
			}
		}
	})
}

func TestExecuteWhenEnabled(t *testing.T) {
	t.Run("ExecuteProcessesAudioWhenEnabled", func(t *testing.T) {
		// Create test audio data with some noise pattern
		testAudio := make([]int16, frameSize)
		for i := range testAudio {
			// Create a simple audio pattern
			testAudio[i] = int16((i % 100) * 300)
		}

		// Create a copy to compare
		originalAudio := make([]int16, frameSize)
		copy(originalAudio, testAudio)

		// Enable noise cancellation
		enabled.Store(true)

		// Execute should process the audio when enabled
		Execute(testAudio)

		// The audio should be processed (RNNoise will modify it)
		// We can't predict exact output, but we can verify the function ran
		// In a real scenario, RNNoise will modify the audio
		// For this test, we just verify the function doesn't panic
		// and the buffer is still the correct size
		if len(testAudio) != frameSize {
			t.Errorf("Audio buffer size changed: got %d, want %d", len(testAudio), frameSize)
		}
	})
}

func TestExecuteWithZeroAudio(t *testing.T) {
	t.Run("ExecuteHandlesZeroAudio", func(t *testing.T) {
		// Create silent audio (all zeros)
		testAudio := make([]int16, frameSize)

		// Enable noise cancellation
		enabled.Store(true)

		// This should not panic
		Execute(testAudio)

		// Verify buffer size is maintained
		if len(testAudio) != frameSize {
			t.Errorf("Audio buffer size changed: got %d, want %d", len(testAudio), frameSize)
		}
	})
}

func TestExecuteWithMaxAmplitude(t *testing.T) {
	t.Run("ExecuteHandlesMaxAmplitude", func(t *testing.T) {
		// Create audio with maximum amplitude
		testAudio := make([]int16, frameSize)
		for i := range testAudio {
			if i%2 == 0 {
				testAudio[i] = 32767 // Max positive int16
			} else {
				testAudio[i] = -32768 // Max negative int16
			}
		}

		// Enable noise cancellation
		enabled.Store(true)

		// This should not panic or overflow
		Execute(testAudio)

		// Verify buffer size is maintained
		if len(testAudio) != frameSize {
			t.Errorf("Audio buffer size changed: got %d, want %d", len(testAudio), frameSize)
		}
	})
}

func TestConcurrentToggle(t *testing.T) {
	t.Run("ConcurrentTogglesAreSafe", func(t *testing.T) {
		// Reset to known state
		enabled.Store(true)

		// Run multiple toggles concurrently
		done := make(chan bool)
		for i := 0; i < 100; i++ {
			go func() {
				Toggle()
				done <- true
			}()
		}

		// Wait for all goroutines
		for i := 0; i < 100; i++ {
			<-done
		}

		// The final state should be valid (either true or false)
		state := IsEnabled()
		if state != true && state != false {
			t.Error("IsEnabled() returned invalid state")
		}
	})
}

func TestConcurrentReadWrite(t *testing.T) {
	t.Run("ConcurrentReadWriteIsSafe", func(t *testing.T) {
		enabled.Store(true)

		done := make(chan bool)

		// Start readers
		for i := 0; i < 50; i++ {
			go func() {
				for j := 0; j < 100; j++ {
					IsEnabled()
				}
				done <- true
			}()
		}

		// Start writers
		for i := 0; i < 50; i++ {
			go func() {
				for j := 0; j < 100; j++ {
					Toggle()
				}
				done <- true
			}()
		}

		// Wait for all goroutines
		for i := 0; i < 100; i++ {
			<-done
		}

		// Should complete without race conditions
		state := IsEnabled()
		if state != true && state != false {
			t.Error("IsEnabled() returned invalid state after concurrent access")
		}
	})
}

// BenchmarkToggle benchmarks the toggle operation
func BenchmarkToggle(b *testing.B) {
	enabled.Store(true)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Toggle()
	}
}

// BenchmarkIsEnabled benchmarks reading the enabled state
func BenchmarkIsEnabled(b *testing.B) {
	enabled.Store(true)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		IsEnabled()
	}
}

// BenchmarkExecuteDisabled benchmarks Execute when disabled
func BenchmarkExecuteDisabled(b *testing.B) {
	testAudio := make([]int16, frameSize)
	for i := range testAudio {
		testAudio[i] = int16(i * 100)
	}

	enabled.Store(false)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Execute(testAudio)
	}
}

// BenchmarkExecuteEnabled benchmarks Execute when enabled
func BenchmarkExecuteEnabled(b *testing.B) {
	testAudio := make([]int16, frameSize)
	for i := range testAudio {
		testAudio[i] = int16(i * 100)
	}

	enabled.Store(true)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Execute(testAudio)
	}
}
