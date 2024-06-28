package noise_canceller

/*
#cgo LDFLAGS: -lrnnoise
#include <stdio.h>
#include <stdlib.h>
#include <rnnoise.h>
*/
import "C"
import "unsafe"

var den *C.DenoiseState

const frameSize = 480

func init() {
	den = C.rnnoise_create(nil)
}

func Execute(inputAudio []int16) {
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

func Close() {
	C.rnnoise_destroy(den)
}
