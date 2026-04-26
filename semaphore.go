package main

import (
	"runtime"
	"sync/atomic"
	"time"
)

type Semaphore struct {
	value int32
}

func Init(initial int32) *Semaphore {
	return &Semaphore{value: initial}
}

func (s *Semaphore) Wait() {
	for {
		v := atomic.LoadInt32(&s.value)
		if v > 0 {
			if atomic.CompareAndSwapInt32(&s.value, v, v-1) {
				return
			}
			continue
		}

		runtime.Gosched()
		time.Sleep(time.Millisecond)
	}
}

func (s *Semaphore) Signal() {
	atomic.AddInt32(&s.value, 1)
}

func (s *Semaphore) GetValue() int32 {
	return atomic.LoadInt32(&s.value)
}
