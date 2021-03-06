package main

import (
	"fmt"
	"bufio"
	"os"
)

func main() {
	running := true
	reader := bufio.NewReader(os.Stdin)
	for running {
		data, _, _ := reader.ReadLine()
		command := string(data)
		if command == "stop" {
			running = false
		}
		fmt.Println("command", command)
	}
}