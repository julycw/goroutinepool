package main

import (
	"fmt"
	"github.com/julycw/goroutinepool"
	"github.com/julycw/goroutinepool/request"
	"runtime"
	"sync"
	"time"
)

var mu sync.Mutex

const (
	min_value = 1
	max_value = 100000
)

func isPrime(number int) bool {
	if number == 1 {
		return false
	} else if number == 2 {
		return true
	}

	for i := 2; i <= number/2; i++ {
		if number%i == 0 {
			return false
		}
	}
	return true
}

func doWork(args ...interface{}) (interface{}, error) {
	a := args[0].(int)
	ret := isPrime(a)
	return ret, nil
}

func withTimestramp(fn func()) {
	begin := time.Now()
	fn()
	eslapse := time.Now().Sub(begin).Nanoseconds() / 10e5
	fmt.Println("Last:", eslapse, "milliseconds")
}

func workBySingleThread() {
	withTimestramp(func() {
		count := 0
		for i := min_value; i <= max_value; i++ {
			if isPrime(i) {
				count++
			}
		}
		fmt.Println("By single thread:", count)
	})
}

func workByGoroutine() {
	withTimestramp(func() {
		goroutines := make(chan int, 10)
		count := 0
		total := 0
		for i := min_value; i <= max_value; i++ {
			goroutines <- 1
			go func(i int, total, count *int) {
				mu.Lock()
				*total++
				mu.Unlock()
				if isPrime(i) {
					mu.Lock()
					*count++
					mu.Unlock()
				}
				<-goroutines
			}(i, &total, &count)
		}
		for {
			time.Sleep(1 * time.Nanosecond)
			if total == max_value-min_value+1 {
				break
			}
		}
		fmt.Println("By Goroutines:", count)
	})
}

func workByPool() {
	withTimestramp(func() {
		pool := goroutinepool.New(5, 4, 10)
		pool.Run()
		go func() {
			for i := min_value; i <= max_value; i = i + 1 {
				work := request.NewWork(doWork, i)
				pool.PutRequest(work, 5*time.Second)
			}
		}()

		total := 0
		count := 0
		for {
			result, err := pool.GetResult(5 * time.Second)
			if err == nil {
				ret := result.Result.(bool)
				total++
				if ret {
					count++
				}
			} else {
				fmt.Println(err.Error())
			}

			if total == max_value-min_value+1 {
				break
			}
		}
		pool.Dismiss()
		fmt.Println("By Goroutine Pool:", count)
	})
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	workBySingleThread()
	workByGoroutine()
	workByPool()
}
