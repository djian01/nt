package mtu

import (
	"fmt"
	"nt/pkg/ntPinger"
	"nt/pkg/ntScan"

	"github.com/spf13/cobra"
)

// Iniital tcpCmd
var mtuCmd = &cobra.Command{
	Use:   "mtu [flags] <URL>", // Sub-command, shown in the -h, Usage field
	Short: "To determine the largest MTU Size to the destination Host/IP",
	Long:  "To determine the largest packet size that can reach a destination without fragmentation.",
	Args:  cobra.ExactArgs(1), // 1 Arg, <Host/IP> is required
	Run:   MtuCommandLink,
	Example: `
** Noted: The global recording functions (-r & -p) are NOT available for the MTU test **

# Example: MTU check for destination google.com"
nt mtu google.com

# Example: MTU check for destination 192.168.1.10 with user defined ceiling test size 9000 set (for Jumbo Frame enabled environment)
nt mtu -s 9000 192.168.1.10
`,
}

// Func - IcmpCommandLink: obtain Flags and call IcmpCommandMain()
func MtuCommandLink(cmd *cobra.Command, args []string) {

	// GFlag -r
	// recording, _ := cmd.Flags().GetBool("recording")

	// GFlag -d
	// displayRow, _ := cmd.Flags().GetInt("displayrow")

	// Flag -s
	ceilingSize, _ := cmd.Flags().GetInt("ceilingsize")

	// Arg - destHost
	destHost := args[0]	

	// call func HttpCommandMain
	err := MtuCommandMain(ceilingSize, destHost)
	if err != nil {
		panic(err) // panic all error from under
	}
}

// Func - HttpCommandMain
func MtuCommandMain(ceilingSize int, destHost string) error {

	// Resolve destHost
	DestAddrs, err := ntPinger.ResolveDestHost(destHost)
	if err != nil {
		return err
	}

	// call the ScanMTURun
	err = ntScan.ScanMTURun(ceilingSize, fmt.Sprint(DestAddrs[0]), destHost)
	if err != nil {
		return err
	}
	return nil
}

// Func - HttpCommand
func MtuCommand() *cobra.Command {

    // Customize the help template for the subcommand. Only the sub command Flag description is shown. Remove the Global Flag description.
	mtuCmd.SetHelpTemplate(`{{.Short}}

Usage:
  {{.UseLine}}

{{if .HasExample}}Example:
{{.Example}}{{end}}

{{if or .HasAvailableLocalFlags .HasAvailablePersistentFlags}}Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}
{{.PersistentFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}
`)
	return mtuCmd
}

// Func - init()
func init() {

	// Flag - Ping count
	var count int
	mtuCmd.Flags().IntVarP(&count, "ceilingsize", "s", 1500, "Ceiling Test Size (default 1500 bytes)")
}
