package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/errakhaoui/noise-canceling/input"
	"github.com/errakhaoui/noise-canceling/noise_canceller"
	"github.com/errakhaoui/noise-canceling/output"
)

func main() {
	log.Println("Start noise cancellation ...")
	log.Println("Press 't' + Enter to toggle noise cancellation ON/OFF")
	log.Printf("Noise cancellation: ENABLED")

	input.StartMicAcquisition()
	defer input.Close()
	defer output.Close()
	defer noise_canceller.Close()

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
