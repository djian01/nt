package root

import (
	"fmt"
	"os"
	"path/filepath"

	"nt/pkg/cmd/ping"

	"github.com/spf13/cobra"
)

// version
var version = "Net-Test version 0.1.1"

// Initial rootCmd
var rootCmd = &cobra.Command{
	Use:   "nt [flags] <sub-command>", // root command, the root command name can or cannot be "greeter" as the executiable file name can change
	Short: "Net-Test CLI",
	Long:  "Net-Test is a set of tools for network testing",
	Run:   RootCommandFunc,
}

// Func - init()
func init() {

	// Flag(s)
	var _report bool
	var _reportPath string
	var _version bool

	//// GFlag - report
	rootCmd.PersistentFlags().BoolVarP(&_report, "report", "r", false, "Enable result-output report")

	//// Get the path of the current executable
	exeFilePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	//// Get the folder containing the executable
	filePath := filepath.Dir(exeFilePath)

	//// GFlag - report path
	rootCmd.PersistentFlags().StringVarP(&_reportPath, "path", "p", filePath, "The output path for result-output report")

	//// GFlag - version
	rootCmd.Flags().BoolVarP(&_version, "version", "v", false, "Show version")

	// Add Sub-Commands
	rootCmd.AddCommand(ping.PingCommand())
}

// Func - RootCommandFunc()
func RootCommandFunc(cmd *cobra.Command, args []string) {
	// Flag -v
	vFlag, _ := cmd.Flags().GetBool("version")

	if vFlag {
		fmt.Printf("%v\n", version) // write output to Stdout for test verification
	}

	//fmt.Fprintf(cmd.OutOrStdout(), "This is a Network Testing Tool Set.\n") // write output to Stdout for test verification
}

// Func - RootCommand()
func RootCommand() *cobra.Command {
	return rootCmd
}
