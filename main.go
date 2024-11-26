package main

import (
	"fmt"
	"os"
)

var version = "unknown"

func main() {
	if err := initRootCommand(version).Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
