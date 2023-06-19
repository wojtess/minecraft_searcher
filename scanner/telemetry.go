package scanner

import chanspeedcounter "minecraft_searcher/chanSpeedCounter"

type Telemetry struct {
	StopChan     <-chan interface{}
	stopChannels []chan interface{}

	controller *WorkersControler

	workerInputDataCounter chanspeedcounter.ChanSpeedCounter[WorkerInputData]
}

func newTelemetry(controller *WorkersControler) Telemetry {
	t := Telemetry{
		StopChan:     make(<-chan interface{}),
		stopChannels: make([]chan interface{}, 0),
		controller:   controller,
	}
	go func() {
		<-t.StopChan
		for _, v := range t.stopChannels {
			v <- 0
		}
	}()
	return t
}

func (t *Telemetry) GetController() *WorkersControler {
	return t.controller
}

func (t *Telemetry) WrapWorkerInputData(input <-chan WorkerInputData) <-chan WorkerInputData {
	out := make(chan WorkerInputData, cap(input))
	counter := chanspeedcounter.NewChanSpeedCounter(input, out, t.addStopChan())
	go counter.Start()
	t.workerInputDataCounter = counter
	return out
}

func (t *Telemetry) GetWorkerInputDataSpeed() int {
	return t.workerInputDataCounter.GetSpeed()
}

func (t *Telemetry) addStopChan() chan interface{} {
	stop := make(chan interface{})
	t.stopChannels = append(t.stopChannels, stop)
	return stop
}
