package ping

import (
	"github.com/spf13/cobra"

	"nt/pkg/sharedStruct"
)

// result[]
var ntResults []sharedStruct.NtResult

// Iniital pingCmd
var pingCmd = &cobra.Command{
	Use:   "ping [flags] <host>", // Sub-command, shown in the -h, Usage field
	Short: "Ping Test Module",
	Long:  "Ping test Module for ICMP testing",
	Args:  cobra.ExactArgs(1), // Only 1 Arg (dest) is required
	Run:   PingCommandFunc,
}

// Func - PingCommandFunc
func PingCommandFunc(cmd *cobra.Command, args []string) {

	// GFlag -p
	path, _ := cmd.Flags().GetString("path")

	// GFlag -r
	report, _ := cmd.Flags().GetBool("report")

	// Arg - dest
	dest := args[0]

	// Flag -c
	count, _ := cmd.Flags().GetInt("count")

	// Flag -s
	size, _ := cmd.Flags().GetInt("size")

	// Flag -i
	interval, _ := cmd.Flags().GetInt("interval")

	// Start Ping Main Command, manually input display Len
	err := ProbingFunc(dest, count, size, interval, report, path, 10)
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
func PingCommand() *cobra.Command {
	return pingCmd
}

// Func - init()
func init() {

	// Flag - Ping count
	var count int
	pingCmd.Flags().IntVarP(&count, "count", "c", 0, "Ping Test Count")

	// Flag - Ping size
	var size int
	pingCmd.Flags().IntVarP(&size, "size", "s", 24, "Ping Test Packet Size (must be larger than 24 Bytes)")

	// Flag - Ping interval
	var interval int
	pingCmd.Flags().IntVarP(&interval, "interval", "i", 1, "Ping Test Interval")
}
