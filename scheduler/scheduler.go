package main

import (
	"errors"
	"sync"
	"time"
)

type Scheduler struct {
	m map[int]*time.Timer

	mutex  sync.Mutex
	closed bool
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		m: make(map[int]*time.Timer),
	}
}

// SetTimeout run some function after some timeout
func (s *Scheduler) SetTimeout(key int, delay time.Duration, action func()) error {
	if action == nil {
		return errors.New("action can't be nil")
	}

	if delay < 0 {
		return errors.New("delay cannot be negative")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.closed {
		return errors.New("scheduler is closed")
	}

	if timer, found := s.m[key]; found {
		timer.Stop()
	}

	s.m[key] = time.AfterFunc(delay, action)
	return nil
}

// CancelTimeout cancel running of some function
func (s *Scheduler) CancelTimeout(key int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if timer, ok := s.m[key]; ok {
		timer.Stop()
		delete(s.m, key)
	}
}

// Close cancel all functions
func (s *Scheduler) Close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.closed = true
	for key, timer := range s.m {
		timer.Stop()
		delete(s.m, key)
	}
}

func main() {

}
