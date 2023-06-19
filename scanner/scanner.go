package scanner

import (
	"fmt"
	"minecraft_searcher/lossesring"
	"os"
)

type scanner struct {
	ordinary      int
	threads       int
	errrorsLogger chan error
	Workers       *WorkersControler
	LatestErrors  lossesring.LossesRing[string]
	Telemetry     Telemetry
}

func RunNewScanner(timeout int, odrinary int, threads int, inputFile *os.File) scanner {
	errorsLogger := make(chan error, threads/2)

	sc := scanner{
		ordinary:      odrinary,
		threads:       threads,
		Workers:       newWorkersControler(errorsLogger, threads, timeout),
		errrorsLogger: errorsLogger,
		LatestErrors:  lossesring.New[string](100),
	}
	sc.Telemetry = newTelemetry(sc.Workers)
	sc.startErrorListener()
	var inputFileReader <-chan WorkerInputData
	inputFileReader = sc.inputFileReader(inputFile)

	inputFileReader = sc.Telemetry.WrapWorkerInputData(inputFileReader)
	workerOutputData := sc.workerOutputParaser()
	sc.Workers.workerOutputData = &workerOutputData
	sc.Workers.workerInputData = &inputFileReader
	sc.Workers.start()
	return sc
}

func (s scanner) startErrorListener() {
	go func() {
		for {
			err := <-s.errrorsLogger
			s.LatestErrors.Push(err.Error())
			// fmt.Printf("Error: %s\n", err)
		}
	}()
}

type WorkersControler struct {
	workers []*Worker
	threads int
	timeout int

	ChangeSize   chan int
	GetWorkers   chan chan []Worker
	SetTimeout   chan int
	GetTimeout   chan chan int
	errorsLogger chan error

	workerInputData  *<-chan WorkerInputData
	workerOutputData *chan<- WorkerOutputData
}

func newWorkersControler(errorsLogger chan error, threads int, timeout int) *WorkersControler {
	changeSize := make(chan int, 100)
	getWorkers := make(chan chan []Worker, 100)

	setTimeout := make(chan int, 10)
	getTimeout := make(chan chan int, 10)

	controller := WorkersControler{
		errorsLogger: errorsLogger,
		ChangeSize:   changeSize,
		GetWorkers:   getWorkers,
		SetTimeout:   setTimeout,
		GetTimeout:   getTimeout,
		threads:      threads,
		timeout:      timeout,
	}
	//func that is handling controller.workers
	go func() {
		for {
			select {
			case size := <-changeSize:
				if size < 0 {
					//decrase
					size = size * -1
					if size > len(controller.workers)-1 {
						size = len(controller.workers) - 1
					}

					i := 0
					for i = 0; i < size; i++ {
						controller.workers[i].Stop()
					}
					controller.workers = controller.workers[i:]
					controller.threads = size
				} else {
					//incrase
					controller.startWorkers(size)
				}
			case v := <-getWorkers:
				out := make([]Worker, 0, len(controller.workers))
				for _, v := range controller.workers {
					out = append(out, *v)
				}
				v <- out
			}
		}
	}()
	//func that is handling controller.timeout
	go func() {
		for {
			select {
			case timeout := <-setTimeout:
				controller.timeout = timeout
			case v := <-getTimeout:
				v <- controller.timeout
			}
		}
	}()
	return &controller
}

func (controller *WorkersControler) start() error {
	//check if we can start
	if controller.workerInputData == nil || controller.workerOutputData == nil {
		return fmt.Errorf("InputDataChannel nad OutputDataChannel is nil")
	}
	//add new workers based on value of "threads" varible
	controller.ChangeSize <- controller.threads
	return nil
}

func (controller *WorkersControler) startWorkers(size int) []*Worker {
	workers := make([]*Worker, 0, size)
	for i := 0; i < size; i++ {
		worker := Worker{
			input:        controller.workerInputData,
			output:       controller.workerOutputData,
			stop:         make(chan interface{}),
			errorsLogger: controller.errorsLogger,
			Id:           i,
			timeout:      controller.timeout,
			State:        Waiting,
		}
		workers = append(workers, &worker)
		worker.start()
	}
	controller.workers = append(controller.workers, workers...)
	return workers
}
