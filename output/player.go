package output

import (
	"github.com/hajimehoshi/oto"
	"log"
)

const (
	sampleRate   = 48000
	bufferSize   = 1024
	channelCount = 1
)

var player *oto.Player

func init() {
	ctx, err := oto.NewContext(sampleRate, channelCount, 2, bufferSize)
	if err != nil {
		log.Fatal(err)
	}
	player = ctx.NewPlayer()
}

func ReadStream(audioStream []int16) {
	b := make([]byte, len(audioStream)*2)
	for i := range audioStream {
		b[2*i] = byte(audioStream[i])
		b[2*i+1] = byte(audioStream[i] >> 8)
	}

	// Play the processed audio
	_, err := player.Write(b)
	if err != nil {
		log.Fatal(err)
	}
}

func Close() {
	if err := player.Close(); err != nil {
		log.Printf("Error closing player: %v", err)
	}
}
