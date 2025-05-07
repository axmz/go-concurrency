package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// SIMPLE
func gen(n int) <-chan func() {
	out := make(chan func())
	go func() {
		for i := range n {
			out <- func() { fmt.Println("# : ", i+1) }
		}
		close(out)
	}()
	return out
}

func simpleWorker(n int) {
	var wg sync.WaitGroup
	tasksCh := make(chan func(), 2)

	go func() {
		for f := range gen(20) {
			tasksCh <- f
		}
		close(tasksCh)
	}()

	wg.Add(n)

	for i := range n {
		go func() {
			defer wg.Done()
			for f := range tasksCh {
				fmt.Println("worker:", i+1, "started task")
				f()
			}
		}()
	}

	wg.Wait()
}

// MEDIUM
type MediumWorkerPool struct {
	wg      sync.WaitGroup
	tasksCh chan func()
}

func (wp *MediumWorkerPool) Wait() {
	wp.wg.Wait()
}

func NewMediumWorkerPool(n int) *MediumWorkerPool {
	w := MediumWorkerPool{
		wg:      sync.WaitGroup{},
		tasksCh: make(chan func(), n),
	}

	w.wg.Add(n)
	for i := range n {
		go func() {
			defer w.wg.Done()
			for f := range w.tasksCh {
				fmt.Println("worker:", i+1, "started task")
				f()
			}
		}()
	}

	return &w
}

func mediumWorker(n int) {
	w := NewMediumWorkerPool(n)

	go func() {
		for f := range gen(20) {
			w.tasksCh <- f
		}
		close(w.tasksCh)
	}()

	w.Wait()
}

// HARD
func hardWorker(n int) {
	w, _ := NewHardWorkerPool(2)

	f := func() {
		fmt.Println("task")
	}

	w.AddTask(f)
	w.AddTask(f)
	w.AddTask(f)
	w.AddTask(f)
}

type HardWorkerPool struct {
	wg      sync.WaitGroup
	taskCh  chan func()
	closeCh chan struct{}

	m      sync.RWMutex
	closed bool
}

func NewHardWorkerPool(n int) (*HardWorkerPool, error) {
	if n <= 0 {
		return nil, errors.New("n should be positive")
	}

	p := &HardWorkerPool{
		taskCh:  make(chan func(), n),
		closeCh: make(chan struct{}),
	}

	go p.startWorkers(n)

	return p, nil
}

func (wp *HardWorkerPool) startWorkers(n int) {
	wp.wg.Add(n)
	for range n {
		go func() {
			defer func() {
				recover() // if we have a problem with f() then we will exit all our n workers.
				wp.wg.Done()
			}()

			for f := range wp.taskCh {
				f()
			}
		}()
	}
	wp.wg.Wait()
	close(wp.closeCh)
}

func (wp *HardWorkerPool) AddTask(task func()) error {
	if task == nil {
		return errors.New("invalid task")
	}

	wp.m.RLock()
	if wp.closed {
		wp.m.RUnlock()
		return errors.New("chan closed")
	}
	wp.m.RUnlock()

	select {
	case wp.taskCh <- task:
		return nil
	default:
		return errors.New("chan blocked")
	}
}

func (wp *HardWorkerPool) Close() {
	if wp.closed {
		return
	}
	wp.m.Lock()
	wp.closed = true
	wp.m.Unlock()

	close(wp.taskCh)
	<-wp.closeCh
}

func main() {
	// simpleWorker(2)
	// mediumWorker(2)
	hardWorker(2)
	time.Sleep(time.Second)
	fmt.Println("exit")
}
