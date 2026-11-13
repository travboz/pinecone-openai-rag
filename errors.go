package main

import (
	"fmt"
	"os"
)

func handleError(err error, msg string) {
	if err != nil {
		fmt.Fprintf(W, "%s: %v\n", msg, err)
		os.Exit(1)
	}
}
