package util

import (
	"fmt"
	"os"
)

func Exit(message string, err error) {
	fmt.Printf(message, err)
	os.Exit(1)
	return
}
