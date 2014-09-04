package goroutinepool

import (
	"errors"
	"github.com/julycw/goroutinepool/request"
	"github.com/julycw/goroutinepool/work_routine"
	"time"
)

type GroutinePoolStatus int

const (
	Initial = iota + 1
	Running
)

type GoroutinePool struct {
	workRoutines []work_routine.WorkGoroutine
	requestQueue work_routine.RequestQueue
	resultQueue  work_routine.ResultQueue
}

func (this *GoroutinePool) AddWorker() {
	new_worker := *work_routine.New(this.requestQueue, this.resultQueue)
	this.workRoutines = append(this.workRoutines, new_worker)
}

func (this *GoroutinePool) AddWorkers(num int) {
	for i := 0; i < num; i++ {
		this.AddWorker()
	}
}

func (this *GoroutinePool) Init(worker_num, request_size, result_size int) *GoroutinePool {
	this.requestQueue = make(chan *request.RequestWork, request_size)
	this.resultQueue = make(chan *request.ResultWork, result_size)
	this.workRoutines = make([]work_routine.WorkGoroutine, 0, worker_num)
	this.AddWorkers(worker_num)
	return this
}

func (this *GoroutinePool) Run() {
	for i, _ := range this.workRoutines {
		go this.workRoutines[i].Run()
	}
}

func (this *GoroutinePool) PutRequest(reque *request.RequestWork, timeout time.Duration) error {
	select {
	case this.requestQueue <- reque:
		return nil
	case <-time.After(timeout):
		return errors.New("Time out!")
	}
}

func (this *GoroutinePool) GetResult(timeout time.Duration) (request.ResultWork, error) {
	select {
	case result := <-this.resultQueue:
		return *result, nil
	case <-time.After(timeout):
		return request.ResultWork{}, errors.New("Time out!")
	}
}

func (this *GoroutinePool) Dismiss() {
	for i, _ := range this.workRoutines {
		this.workRoutines[i].Stop()
	}
}

func New(worker_num, request_size, result_size int) *GoroutinePool {
	return new(GoroutinePool).Init(worker_num, request_size, result_size)
}
