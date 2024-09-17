package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/djian01/nt/pkg/cmd/root" // import root pkg
)

// create a global logFile pointer and logger pointer
var (
	logFile *os.File
	logger  *log.Logger
)

func main() {

	// create or open the output.txt file for logging
	// "os.O_RDWR": open file to read and write
	// "os.O_CREATE": Create the file with the mode permissions if file does not exist. Cursor is at the beginning.
	// "os.O_APPEND": Only allow write past end of file
	logFile, err := os.OpenFile("nt.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file: ", err)
		return
	}
	defer logFile.Close()

	// create a new logger
	logger = log.New(logFile, "", log.LstdFlags)

	//// defer func() to capture the panic & debug stack messages
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			logger.Printf("Recovered panic: %v", r)
			stack := debug.Stack()
			logger.Printf("Stack Trace: %v", string(stack))
		}
	}()

	// call the rootCmd
	rootCmd := root.RootCommand()

	err = rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
