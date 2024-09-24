package main

import (
	"fmt"
	"time"
)

// there is more examples in the book
func main() {
	doWork := func(
		done <-chan interface{},
		pulseInterval time.Duration,
	) (<-chan interface{}, <-chan time.Time) {
		heartbeat := make(chan interface{})
		results := make(chan time.Time)

		go func() {
			defer close(heartbeat)
			defer close(results)

			pulse := time.Tick(pulseInterval)
			workGen := time.Tick(2 * pulseInterval)

			sendPulse := func() {
				select {
				case heartbeat <- struct{}{}:
				default:
				}
			}

			sendResult := func(r time.Time) {
				for {
					select {
					// The program can exit immediately if done is triggered (nested done case), avoiding a situation where it blocks while trying to send a result.
					case <-done:
						return
					// The program continues sending heartbeat pulses while waiting to send the result.
					case <-pulse:
						sendPulse()
					case results <- r:
						return
					}
				}
			}

			for {
				select {
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case r := <-workGen:
					sendResult(r)
				}
			}
		}()
		return heartbeat, results
	}

	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() { close(done) })
	const timeout = 2 * time.Second
	heartbeat, results := doWork(done, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			if !ok {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if !ok {
				return
			}
			fmt.Printf("results %v\n", r.Second())
		case <-time.After(timeout):
			return
		}
	}
}

// incorrect example
// func doWork2(done <-chan struct{}, interval time.Duration) (chan time.Time, chan time.Time) {
// 	heartbeat := make(chan time.Time, 1)
// 	result := make(chan time.Time, 1)

// 	pulse := time.Tick(interval)
// 	work := time.Tick(interval * 2)

// 	go func() {
// 		defer close(heartbeat)
// 		defer close(result)

// 		for {
// 			// This makes a big difference. Why?
// 			// select {
// 			// case heartbeat <- <-pulse:
// 			// 	fmt.Println("send pulse")
// 			// case result <- <-work:
// 			// 	fmt.Println("send work")
// 			// case <-done:
// 			// 	return
// 			// }

// 			select {
// 			case h := <-pulse:
// 				heartbeat <- h
// 				fmt.Println("send pulse")
// 			case r := <-work:
// 				result <- r
// 				fmt.Println("send work")
// 			case <-done:
// 				return
// 			}
// 		}
// 	}()

// 	return heartbeat, result
// }

// func main() {
// 	done := make(chan struct{})
// 	heartbeat, work := doWork2(done, time.Second)
// 	timeout := time.Second * 3

// 	go func() {
// 		<-time.After(time.Second * 10)
// 		close(done)
// 	}()

// loop:
// 	for {
// 		select {
// 		case <-done:
// 			fmt.Println("done")
// 			break loop
// 		case h := <-heartbeat:
// 			fmt.Println("pulse: ", h)
// 		case w := <-work:
// 			fmt.Println("work: ", w)
// 		case <-time.After(timeout):
// 			fmt.Println("timeout")
// 			break loop
// 		}
// 	}

// 	fmt.Println("exit")
// }
