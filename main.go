package main

import (
	"os"

	"github.com/zot24/repowiki/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
