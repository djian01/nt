package root

import (
	"fmt"
	"os"
	"path/filepath"

	"nt/pkg/cmd/icmp"

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

// Func - RootCommandFunc()
func RootCommandFunc(cmd *cobra.Command, args []string) {
	// Flag -v
	vFlag, _ := cmd.Flags().GetBool("version")

	if vFlag {
		fmt.Printf("%v\n", version) // write output to Stdout for test verification
	}
}

// Func - RootCommand()
func RootCommand() *cobra.Command {
	return rootCmd
}

// Func - init()
func init() {
	// Flag(s)
	var _recording bool
	var _reportPath string
	var _version bool
	var _displayRow int

	//// GFlag - report
	rootCmd.PersistentFlags().BoolVarP(&_recording, "recording", "r", false, "Enable result recording to CSV file")

	//// Get the path of the current executable
	exeFilePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	//// Get the folder containing the executable
	filePath := filepath.Dir(exeFilePath)

	//// GFlag - recording path
	rootCmd.PersistentFlags().StringVarP(&_reportPath, "path", "p", filePath, "The output path for result-output report")

	//// GFlag - display Row Length
	rootCmd.PersistentFlags().IntVarP(&_displayRow, "displayrow", "d", 10, "Set the length of the dispaly row")

	//// Flag - version
	rootCmd.Flags().BoolVarP(&_version, "version", "v", false, "Show version")

	// Add Sub-Commands
	rootCmd.AddCommand(icmp.IcmpCommand())
}
