package chanspeedcounter

import "time"

type ChanSpeedCounter[T any] struct {
	input  <-chan T
	output chan<- T
	stop   chan interface{}
	speed  []int64
}

func NewChanSpeedCounter[T any](input <-chan T, output chan<- T, stop chan interface{}) ChanSpeedCounter[T] {
	return ChanSpeedCounter[T]{
		input:  input,
		output: output,
		stop:   stop,
		speed:  make([]int64, 0),
	}
}

func (counter *ChanSpeedCounter[T]) Start() {
	go func() {
	loop:
		for {
			i := 0
			for _, v := range counter.speed {
				if !(v < time.Now().UnixMilli()-1000) {
					counter.speed[i] = v
					i++
				}
			}
			counter.speed = counter.speed[:i]

			// for i, v := range counter.speed {
			// 	if v < time.Now().UnixMilli()-1000 {
			// 		//remove data that is older than 1sec
			// 		counter.speed[i] = counter.speed[len(counter.speed)-1]
			// 		counter.speed = counter.speed[:len(counter.speed)-1]
			// 	}
			// }
			select {
			case data := <-counter.input:
				counter.speed = append(counter.speed, time.Now().UnixMilli())
				counter.output <- data
			case <-counter.stop:
				break loop
			}
		}
	}()
}

func (counter *ChanSpeedCounter[_]) GetSpeed() int {
	return len(counter.speed)
}
