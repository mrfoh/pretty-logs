package main

import (
	"github.com/mrfoh/pretty-logs/cmd/prettylogs"
)

func main() {
	// Root command
	rootCmd := prettylogs.NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
