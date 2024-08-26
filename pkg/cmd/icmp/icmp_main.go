package icmp

import (
	"github.com/spf13/cobra"

	"nt/pkg/sharedStruct"
)

// result[]
var ntResults []sharedStruct.NtResult

// Iniital pingCmd
var icmpCmd = &cobra.Command{
	Use:   "ping [flags] <host>", // Sub-command, shown in the -h, Usage field
	Short: "Ping Test Module",
	Long:  "Ping test Module for ICMP testing",
	Args:  cobra.ExactArgs(1), // Only 1 Arg (dest) is required
	Run:   IcmpCommandFunc,
}

// Func - PingCommandFunc: the linkage between cobra.Command and the Probing func
func IcmpCommandFunc(cmd *cobra.Command, args []string) {

	// GFlag -p
	path, _ := cmd.Flags().GetString("path")

	// GFlag -r
	report, _ := cmd.Flags().GetBool("report")

	// GFlag -d
	displayRow, _ := cmd.Flags().GetInt("displayrow")

	// Arg - dest
	dest := args[0]

	// Flag -c
	count, _ := cmd.Flags().GetInt("count")

	// Flag -s
	size, _ := cmd.Flags().GetInt("size")

	// Flag -i
	interval, _ := cmd.Flags().GetInt("interval")

	// Start Ping Main Command, manually input display Len
	err := IcmpProbingFunc(dest, count, size, interval, report, path, displayRow)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("ping: %v\n", dest)
	// fmt.Printf("ping count: %v\n", count)

	// fmt.Printf("GFlag path: %v\n", path)
	// fmt.Printf("GFlag report: %v\n", report)

	// fmt.Println(ntResults)
}

// Func - PingCommand
func IcmpCommand() *cobra.Command {
	return icmpCmd
}

// Func - init()
func init() {

	// Flag - Ping count
	var count int
	icmpCmd.Flags().IntVarP(&count, "count", "c", 0, "Ping Test Count")

	// Flag - Ping size
	var size int
	icmpCmd.Flags().IntVarP(&size, "size", "s", 24, "Ping Test Packet Size (must be larger than 24 Bytes)")

	// Flag - Ping interval
	var interval int
	icmpCmd.Flags().IntVarP(&interval, "interval", "i", 1, "Ping Test Interval")
}
