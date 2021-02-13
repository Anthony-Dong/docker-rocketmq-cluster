package common

import (
	"fmt"
	"os"
	"time"
)

func EchoError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "find error: %v\n", err)
	}
}

func NowTimeString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}