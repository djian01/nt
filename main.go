package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"

	"github.com/djian01/nt/pkg/cmd/root" // import root pkg
)

// create a global logFile pointer and logger pointer
var (
	logFile *os.File
	logger  *log.Logger
)

// Func: get the config file path for different OS
func getConfigFilePath(appName string) (string, error) {
	var configDir string
	var err error

	if runtime.GOOS == "darwin" {
		// macOS: ~/Library/Application Support/<appName> (/Users/<User Name>/Library/Application Support/<appName>)
		configDir, err = os.UserConfigDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(configDir, appName)
	} else {
		// Windows/Linux: directory where executable resides
		exePath, err := os.Executable()
		if err != nil {
			return "", err
		}
		configDir = filepath.Dir(exePath)
	}

	// Ensure the config directory exists
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		return "", err
	}

	// Return full path for config file
	return configDir, nil
}

func main() {

	// get the config file path
	// macOS: ~/Library/Application Support/<appName>
	// Windows & Linux: the config file path is the same as the executable path
	configPath, err := getConfigFilePath("nt")
	if err != nil {
		log.Fatal("Failed to get log file path:", err)
		return
	}

	// create or open the output.txt file for logging
	// "os.O_RDWR": open file to read and write
	// "os.O_CREATE": Create the file with the mode permissions if file does not exist. Cursor is at the beginning.
	// "os.O_APPEND": Only allow write past end of file
	logFile, err := os.OpenFile(filepath.Join(configPath, "nt.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
