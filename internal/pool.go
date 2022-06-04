package internal

import "fmt"

type Task interface {
	Execute()
}

type Executor struct {
	Routines uint
	Tasks chan Task
}

func NewExecutor(parallelism uint) *Executor {
	executor := Executor{
		Routines: parallelism,
		Tasks:    make(chan Task),
	}

	for i := uint(0); i < executor.Routines; i++ {
		go executor.lifecycle()
	}

	return &executor
}

func (ex *Executor) Queue(t Task) {
	fmt.Println("Enqueue job\n")
	ex.Tasks <- t
}

func (ex *Executor) lifecycle() {
	for {
		select {
		case task := <-ex.Tasks:
			task.Execute()
			break
		}
	}
}