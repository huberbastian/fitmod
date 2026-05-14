package main

import (
	"os"

	"github.com/huberbastian/fitmod/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
