package main

import (
	"fmt"
	"os"

	"github.com/mj0nez/restic-exporter/cmd"
)

func main() {

	os.Exit(func() int {
		// fmt.Printf("Arguments %v\n", os.Args[1:])

		if err := cmd.Execute(); err != nil {
			fmt.Fprintf(os.Stderr, "Encountered error %v\n", err)
			return 1
		} else {
			return 0
		}

	}())
}
