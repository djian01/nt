package main

import (
	"fmt"
	"nt/pkg/cmd/root" // import root pkg
	"os"
)

func main() {
	rootCmd := root.RootCommand()

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
