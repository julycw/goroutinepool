package work_routine

import (
	"errors"
	"fmt"
	"github.com/julycw/goroutinepool/request"
	"time"
)

type GoRoutineStatus int

type RequestQueue chan *request.RequestWork
type ResultQueue chan *request.ResultWork

const (
	Waiting GoRoutineStatus = iota + 1
	Busy
	Stopped
	Stopping
)

func (this GoRoutineStatus) String() string {
	switch this {
	case Waiting:
		return "Waiting"
	case Busy:
		return "Busy"
	case Stopped:
		return "Stopped"
	case Stopping:
		return "Stopping"
	default:
		return "Unknow"
	}
}

type WorkGoroutine struct {
	requestQueue RequestQueue
	resultQueue  ResultQueue
	status       GoRoutineStatus
}

func (this *WorkGoroutine) Init(requestQueue RequestQueue, resultQueue ResultQueue) *WorkGoroutine {
	this.requestQueue = requestQueue
	this.resultQueue = resultQueue
	this.status = Waiting
	return this
}

func (this *WorkGoroutine) Run() {
	this.status = Waiting
	for {
		if this.status == Stopping {
			this.status = Stopped
			break
		}

		if work, err := this.GetWork(5 * time.Second); err == nil {
			this.status = Busy
			ret, err := work.Call()
			if err != nil {
				work.HasError = true
			}
			this.resultQueue <- request.NewResult(work, ret, err)
			this.status = Waiting
		} else {
			fmt.Println(err.Error())
		}
	}
}

func (this *WorkGoroutine) Stop() {
	if this.status == Busy || this.status == Waiting {
		this.status = Stopping
	}
}

func (this *WorkGoroutine) IsBusy() bool {
	return this.status == Busy
}

func (this *WorkGoroutine) IsStop() bool {
	return this.status == Stopped
}

func (this *WorkGoroutine) GetWork(timeout time.Duration) (*request.RequestWork, error) {
	select {
	case ret := <-this.requestQueue:
		return ret, nil
	case <-time.After(timeout):
		return nil, errors.New("No more work!")
	}
}

func New(requestQueue RequestQueue, resultQueue ResultQueue) *WorkGoroutine {
	return new(WorkGoroutine).Init(requestQueue, resultQueue)
}
