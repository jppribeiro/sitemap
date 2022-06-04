package internal

import (
	"fmt"
	"sync/atomic"
)

type Task interface {
	Execute(r uint)
}

type Executor struct {
	Routines uint
	Tasks chan Task
	Capacity uint32
}

func NewExecutor(parallelism uint) *Executor {
	executor := Executor{
		Routines: parallelism,
		Tasks:    make(chan Task, 10000),
		Capacity: 500,
	}

	for i := uint(0); i < executor.Routines; i++ {
		go executor.lifecycle(i)
	}

	return &executor
}

func (ex *Executor) Queue(t Task) error {
	if atomic.AddUint32(&ex.Capacity, 0) == 0 {
		return fmt.Errorf("Queue is full")
	}

	ex.Tasks <- t

	fmt.Println("entuqued")
	atomic.AddUint32(&ex.Capacity, ^uint32(0))
	return nil
}

func (ex *Executor) lifecycle(r uint) {
	for {
		select {
		case task := <-ex.Tasks:
			atomic.AddUint32(&ex.Capacity, ^uint32(0))
			task.Execute(r)
			atomic.AddUint32(&ex.Capacity, 1)
			break
		}
	}
}