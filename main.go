package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		os.Exit(1)
	}

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		os.Exit(1)
	}

	fmt.Println("starting server on port ", port)

	RunMultiThreadedServer(port)
}
