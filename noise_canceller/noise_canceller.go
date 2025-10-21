package noise_canceller

/*
#cgo LDFLAGS: -lrnnoise
#include <stdio.h>
#include <stdlib.h>
#include <rnnoise.h>
*/
import "C"
import (
	"sync/atomic"
	"unsafe"
)

var den *C.DenoiseState
var enabled atomic.Bool

const frameSize = 480

func init() {
	den = C.rnnoise_create(nil)
	enabled.Store(true) // Start with noise cancellation enabled
}

func Execute(inputAudio []int16) {
	// Only process if noise cancellation is enabled
	if !enabled.Load() {
		return
	}

	inFloat := make([]C.float, frameSize)
	for i := range inputAudio {
		inFloat[i] = C.float(inputAudio[i])
	}

	// Apply RNNoise to the audio frame
	C.rnnoise_process_frame(den, (*C.float)(unsafe.Pointer(&inFloat[0])), (*C.float)(unsafe.Pointer(&inFloat[0])))

	// Convert float samples back to int16
	for i := range inFloat {
		inputAudio[i] = int16(inFloat[i])
	}
}

// Toggle switches noise cancellation on/off
func Toggle() bool {
	newState := !enabled.Load()
	enabled.Store(newState)
	return newState
}

// IsEnabled returns the current state of noise cancellation
func IsEnabled() bool {
	return enabled.Load()
}

func Close() {
	C.rnnoise_destroy(den)
}
