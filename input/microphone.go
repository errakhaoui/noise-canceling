package input

import (
	"github.com/gordonklaus/portaudio"
	"log"
)

const SampleRate = 48000
const frameSize = 480

var InputBuffer []int16
var inputStream *portaudio.Stream

func init() {
	if err := portaudio.Initialize(); err != nil {
		log.Fatal(err)
	}
}

func StartMicAcquisition() {
	var err error
	InputBuffer = make([]int16, frameSize)

	inputStream, err = portaudio.OpenDefaultStream(1, 0, SampleRate, len(InputBuffer), InputBuffer)
	if err != nil {
		log.Fatal(err)
	}

	err = inputStream.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func ReadStream() {
	err := inputStream.Read()
	if err != nil {
		log.Fatal(err)
	}
}

func Close() {
	if err := inputStream.Stop(); err != nil {
		log.Printf("Error stopping input stream: %v", err)
	}
	if err := inputStream.Close(); err != nil {
		log.Printf("Error closing input stream: %v", err)
	}
	if err := portaudio.Terminate(); err != nil {
		log.Printf("Error terminating portaudio: %v", err)
	}
}
