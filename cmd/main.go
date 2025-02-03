package main

import (
	"github.com/mrfoh/pretty-logs/cmd/prettylogs"
)

func main() {
	// Root command
	rootCmd := prettylogs.NewRootCmd()

	rootCmd.AddCommand(prettylogs.NewVersionCmd())

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
