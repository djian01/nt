package root

import (
	"fmt"

	"github.com/djian01/nt/pkg/cmd/dns"
	Http "github.com/djian01/nt/pkg/cmd/http"
	"github.com/djian01/nt/pkg/cmd/icmp"
	"github.com/djian01/nt/pkg/cmd/mtu"
	"github.com/djian01/nt/pkg/cmd/tcp"

	"github.com/spf13/cobra"
)

// version
var version = "0.3.5"

// Initial rootCmd
var rootCmd = &cobra.Command{
	Use:   "nt [flags] <sub-command: icmp/tcp/http/dns/mtu>", // root command, the root command name can or cannot be "greeter" as the executiable file name can change
	Short: "Net-Test CLI",
	Long:  "Net-Test is a set of tools for network testing",
	Run:   RootCommandFunc,
	Example: `
# Example: ICMP ping to "google.com" with recording enabled
nt -r icmp google.com

# Example: TCP ping to "10.2.3.10:22" with count: 10 and interval: 2 sec
nt tcp -c 10 -i 2 10.2.3.10 22
`,
}

// Func - RootCommandFunc()
func RootCommandFunc(cmd *cobra.Command, args []string) {
	// Flag -v
	vFlag, _ := cmd.Flags().GetBool("version")

	if vFlag {
		fmt.Println("")
		fmt.Printf("Net-Test Version:   %v\n", version) // write output to Stdout for test verification
		fmt.Printf("%v\n", "Developed By:       Dennis Jian")
		fmt.Printf("%v\n", "Project Home:       https://github.com/djian01/nt")
		fmt.Println("")
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
	rootCmd.PersistentFlags().IntVarP(&_displayRow, "displayrow", "p", 10, "Set the number of the dispaly row(s)")

	//// Flag - version
	rootCmd.Flags().BoolVarP(&_version, "version", "v", false, "Show version")

	// Add Sub-Commands
	rootCmd.AddCommand(tcp.TcpCommand())
	rootCmd.AddCommand(icmp.IcmpCommand())
	rootCmd.AddCommand(Http.HttpCommand())
	rootCmd.AddCommand(dns.DnsCommand())
	rootCmd.AddCommand(mtu.MtuCommand())
}
