package main

import (
	"context"
	"fmt"
	"time"
)

func shutdown1() {
	stopChan := make(chan struct{})
	shutdownCh := make(chan struct{})

	go func() {
		<-stopChan
		fmt.Println("Shutting down gracefully")
		time.Sleep(2 * time.Second)
		close(shutdownCh)
	}()
	// Simulate work
	time.Sleep(2 * time.Second)
	close(stopChan)
	<-shutdownCh
	fmt.Println("exit")
}

func shutdown2() {
	ctx, cancel := context.WithCancel(context.Background())
	shutdownCh := make(chan struct{})

	go func() {
		<-ctx.Done()
		fmt.Println("Shutting down gracefully")
		time.Sleep(2 * time.Second)
		close(shutdownCh)
	}()

	// Simulate work
	time.Sleep(2 * time.Second)
	cancel()
	select {
	case <-shutdownCh:
		fmt.Println("Shutdown completed")
	case <-time.After(3 * time.Second):
		fmt.Println("Shutdown timed out")
	}
	fmt.Println("exit")
}

func main() {
	shutdown2()
}
