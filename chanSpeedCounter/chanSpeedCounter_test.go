package chanspeedcounter_test

import (
	chanspeedcounter "minecraft_searcher/chanSpeedCounter"
	"testing"
	"time"
)

func TestChanSpeedCounter(t *testing.T) {
	output := make(chan int, 10)

	//read all data
	go func() {
		for {
			<-output
		}
	}()

	input := make(chan int, 10)

	go func() {
		for {
			input <- 1
		}
	}()

	speedCounter := chanspeedcounter.NewChanSpeedCounter(input, output, make(chan interface{}))
	speedCounter.Start()

	for {
		time.Sleep(time.Millisecond * 500)
		t.Logf("Speed: %dps\n", speedCounter.GetSpeed())
	}

}
