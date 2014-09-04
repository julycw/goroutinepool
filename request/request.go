package request

import (
	"errors"
	"math/rand"
	"time"
)

type ResultWork struct {
	Request *RequestWork
	Result  interface{}
	Error   error
}

func (this *ResultWork) Init(request *RequestWork, result interface{}, err error) *ResultWork {
	this.Request = request
	this.Result = result
	this.Error = err
	return this
}

type RequestWork struct {
	Id       uint64
	HasError bool
	Args     []interface{}
	callable func(args ...interface{}) (interface{}, error)
}

func NewResult(request *RequestWork, result interface{}, err error) *ResultWork {
	return new(ResultWork).Init(request, result, err)
}

func (this *RequestWork) Init(callable func(args ...interface{}) (interface{}, error), args ...interface{}) *RequestWork {
	this.Id = uint64(time.Now().Unix()*100000) + uint64(rand.Int63n(100000))
	this.callable = callable
	this.Args = args
	this.HasError = false
	return this
}

func NewWork(callable func(args ...interface{}) (interface{}, error), args ...interface{}) *RequestWork {
	return new(RequestWork).Init(callable, args...)
}

func (this *RequestWork) Call() (interface{}, error) {
	if this.callable == nil {
		return nil, errors.New("Callable is nil.")
	}
	return this.callable(this.Args...)
}

func init() {
	rand.Seed(time.Now().Unix())
}
