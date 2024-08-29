package root

import (
	"fmt"

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
	var _version bool
	var _displayRow int

	//// GFlag - report
	rootCmd.PersistentFlags().BoolVarP(&_recording, "recording", "r", false, "Enable result recording to a CSV file")

	//// GFlag - display Row Length
	rootCmd.PersistentFlags().IntVarP(&_displayRow, "displayrow", "d", 10, "Set the number of the dispaly row(s)")

	//// Flag - version
	rootCmd.Flags().BoolVarP(&_version, "version", "v", false, "Show version")

	// Add Sub-Commands
	rootCmd.AddCommand(icmp.IcmpCommand())
}
