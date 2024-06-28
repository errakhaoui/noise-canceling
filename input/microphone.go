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
	portaudio.Initialize()
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
	portaudio.Terminate()
	inputStream.Stop()
	inputStream.Close()
}
