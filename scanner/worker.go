package scanner

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Tnze/go-mc/chat"
	"github.com/google/uuid"
)

type State int

const (
	Waiting  State = 0
	Working  State = 1
	Stopping State = 2
)

type Worker struct {
	input        *<-chan WorkerInputData
	stop         chan interface{}
	output       *chan<- WorkerOutputData
	errorsLogger chan error
	Id           int
	timeout      int
	State        State
	LastOutput   WorkerOutputData
	LastInput    WorkerInputData
}

type WorkerOutputData struct {
	Status     status
	IdOfWorker int
	Ip         string
	Delay      time.Duration
}

type WorkerInputData struct {
	Ip       string
	Ordinary int
}

type status struct {
	Description chat.Message
	Players     struct {
		Max    int
		Online int
		Sample []struct {
			ID   uuid.UUID
			Name string
		}
	}
	Version struct {
		Name     string
		Protocol int
	}
	Favicon string
	Delay   int64
}

func (w *Worker) Stop() {
	w.State = Stopping
	w.stop <- 0
}

func (w *Worker) start() {
	w.State = Working
	go func() {
	loop:
		for {
			select {
			case input := <-*w.input:
				// fmt.Printf("[%d]ip: %s ordinaty: %d\n", w.id, input.ip, input.ordinary)
				w.LastInput = input
				if input.Ordinary == -1 {
					break loop
				}
				for {
					ip := input.Ip
					resp, delay, err := pingServer(ip, time.Duration(w.timeout)*time.Millisecond)
					if err != nil {
						if strings.HasSuffix(err.Error(), "socket: too many open files") {
							time.Sleep(time.Duration(100) * time.Millisecond)
							continue
						} else {
							w.errorsLogger <- fmt.Errorf("worker %d encountered error: %s", w.Id, err)
							break
						}
					}
					// fmt.Printf("[%d]found \"%s\" after %dms \n", input.ordinary, ip, delay.Milliseconds())
					var s status
					err = json.Unmarshal([]byte(resp), &s)
					if err != nil {
						w.errorsLogger <- fmt.Errorf("worker %d encountered error: %s", w.Id, err)
					}
					s.Delay = delay.Milliseconds()
					w.LastOutput = WorkerOutputData{
						Status:     s,
						IdOfWorker: w.Id,
						Ip:         ip,
						Delay:      delay,
					}
					*w.output <- w.LastOutput
					break
				}
			case <-w.stop:
				break loop
			}
		}
	}()
}
