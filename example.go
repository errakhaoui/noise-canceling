package main

import (
	"github.com/errakhaoui/noise-canceling/input"
	"github.com/errakhaoui/noise-canceling/noise_canceller"
	"github.com/errakhaoui/noise-canceling/output"
	"log"
)

func main() {
	log.Println("Start noise cancellation ...")
	input.StartMicAcquisition()
	defer input.Close()
	defer output.Close()
	defer noise_canceller.Close()

	for {
		// Read audio from the input stream
		input.ReadStream()
		noise_canceller.Execute(input.InputBuffer)
		output.ReadStream(input.InputBuffer)
	}
}
