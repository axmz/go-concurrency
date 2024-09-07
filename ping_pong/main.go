package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	ping := make(chan string)
	pong := make(chan string)

	go func(ping, pong chan string) {
		for {
			str := <-ping
			time.Sleep(1 * time.Second)
			pong <- strings.ToUpper(str)
		}

	}(ping, pong)

	for {
		fmt.Print("Enter ping: ")
		var input string
		fmt.Scanln(&input)
		if input == "q" {
			close(ping)
			close(pong)
			break
		}

		ping <- input
		fmt.Println(<-pong)
	}

	fmt.Println("Exited")
}
