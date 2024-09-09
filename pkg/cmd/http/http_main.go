package Http

import (
	"fmt"
	"net/url"
	"nt/pkg/ntPinger"
	"nt/pkg/record"
	"nt/pkg/terminalOutput"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Iniital tcpCmd
var httpCmd = &cobra.Command{
	Use:   "http [flags] <URL>", // Sub-command, shown in the -h, Usage field
	Short: "http Ping Test Module",
	Long:  "http Ping test Module for http testing",
	Args:  cobra.ExactArgs(1), // 1 Args, <url> are required
	Run:   HttpCommandLink,
	Example: `
# Example: HTTP ping to "https://google.com" with recording enabled. Default Values: Port-443, Method-GET, Count-0, Interval-5s, Timeout-4s
nt -r http https://google.com

# Example: HTTP ping to POST "http://10.2.3.10:8080/token" with count: 10 and interval: 2 sec
nt http -c 10 -i 2 -m POST http://10.2.3.10:8080/token
`,
}

// Initial the bucket
var bucket = 10

// Func - IcmpCommandLink: obtain Flags and call IcmpCommandMain()
func HttpCommandLink(cmd *cobra.Command, args []string) {

	// GFlag -r
	recording, _ := cmd.Flags().GetBool("recording")

	// GFlag -d
	displayRow, _ := cmd.Flags().GetInt("displayrow")

	// Arg - HttpVarInput
	HttpVarInput, err := ParseURL(args[0])
	if err != nil {
		panic(err)
	}

	// Flag -c
	count, _ := cmd.Flags().GetInt("count")

	// Flag -t
	timeout, _ := cmd.Flags().GetInt("timeout")

	// Flag -i
	interval, _ := cmd.Flags().GetInt("interval")

	// Flag -m
	HttpMethod, _ := cmd.Flags().GetString("method")

	// call func HttpCommandMain
	err = HttpCommandMain(recording, displayRow, HttpVarInput, HttpMethod, count, timeout, interval)
	if err != nil {
		panic(err) // panic all error from under
	}
}

// Func - HttpCommandMain
func HttpCommandMain(recording bool, displayRow int, HttpVarInput HttpVar, HttpMethod string, count int, timeout int, interval int) error {

	// Wait Group
	var wgRecord sync.WaitGroup

	// recording row
	recordingRow := 0
	if recording {
		recordingRow = 1
	}

	// recordingFilePath
	recordingFilePath := ""

	// Channel - outputChan (if there are N go routine, the channel deep is N)
	outputChan := make(chan ntPinger.Packet, 1)
	defer close(outputChan)

	// Channel - error (for Go Routines)
	errChan := make(chan error, 1)
	defer close(errChan)

	// Channel - recordingChan, closed in the end of the testing, no need to defer close
	recordingChan := make(chan ntPinger.Packet, 1)

	// build the InputVar
	InputVar := ntPinger.InputVars{
		Type:        "http",
		Count:       count,
		Timeout:     timeout,
		Interval:    interval,
		DestHost:    HttpVarInput.Hostname,
		DestPort:    HttpVarInput.Port,
		Http_scheme: HttpVarInput.Scheme,
		Http_method: HttpMethod,
		Http_path:   HttpVarInput.Path,
	}

	// Start Ping Main Command, manually input display Len
	p, err := ntPinger.NewPinger(InputVar)
	if err != nil {
		return err // return err from NewPinger including resolve error
	}

	go p.Run(errChan)

	// Output
	//// Go Routine: OutputFunc
	go terminalOutput.OutputFunc(outputChan, displayRow, recording)

	// Recording
	if recording {

		// recordingFile Path
		exeFileFolder, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		// recordingFile Name
		timeStamp := time.Now().Format("20060102150405")
		recordingFileName := fmt.Sprintf("Record_%v_%v_%v.csv", "tcp", HttpVarInput.Hostname, timeStamp)
		recordingFilePath = filepath.Join(exeFileFolder, recordingFileName)

		// Go Routine: RecordingFunc
		go record.RecordingFunc(recordingFilePath, bucket, recordingChan, &wgRecord)
	}

	// harvest the result
	loopClose := false
	for {
		// check loopClose Flag
		if loopClose {
			break
		}

		// select option
		select {
		case pkt, ok := <-p.ProbeChan:
			if !ok {
				loopClose = true
				break // break select, bypass "outputChan <- pkt"
			}

			// outputChan
			outputChan <- pkt

			// recordingChan
			if recording {
				recordingChan <- pkt
			}
		case err := <-errChan:
			return err
		}
	}

	// wait for the last interval
	time.Sleep(time.Duration(1) * time.Second)

	// if recording Enabled
	if recording {
		wgRecord.Add(1)
		close(recordingChan)
		// waiting the recording function to save the last records
		wgRecord.Wait()
	}

	// display testing completed
	fmt.Printf("\033[%d;1H", (displayRow + recordingRow + 7))
	fmt.Println("\n--- testing completed ---")

	// if recording is enabled, display the recording file path
	if recording {
		fmt.Printf("Recording CSV file is saved at: %s\n", color.GreenString(recordingFilePath))
	}

	return nil
}

// Func - IcmpCommand
func HttpCommand() *cobra.Command {
	return httpCmd
}

// Func - init()
func init() {

	// Flag - Ping count
	var count int
	httpCmd.Flags().IntVarP(&count, "count", "c", 0, "HTTP Ping Count (default 0 - Non Stop till Ctrl+C)")

	// Flag - HTTP Method
	var method string
	httpCmd.Flags().StringVarP(&method, "method", "m", "GET", "HTTP Ping Metohd (default: GET)")

	// Flag - Ping timeout
	var timeout int
	httpCmd.Flags().IntVarP(&timeout, "timeout", "t", 4, "HTTP Ping Timeout (default: 4 sec)")

	// Flag - Ping interval
	var interval int
	httpCmd.Flags().IntVarP(&interval, "interval", "i", 5, "HTTP Ping Interval (default: 5 sec)")
}

type HttpVar struct {
	Scheme   string
	Hostname string
	Port     int
	Path     string
}

// ParseURL extracts scheme, hostname, port, and path from a URL
func ParseURL(inputURL string) (HttpVar, error) {

	HttpVarNew := HttpVar{}

	parsedURL, err := url.Parse(inputURL)

	if err != nil {
		return HttpVarNew, err
	}

	HttpVarNew.Scheme = parsedURL.Scheme
	HttpVarNew.Hostname = parsedURL.Hostname()

	// Handle default ports for http and https
	if parsedURL.Port() != "" {
		HttpVarNew.Port, err = strconv.Atoi(parsedURL.Port())
		if err != nil {
			return HttpVarNew, err
		}
	} else if HttpVarNew.Scheme == "http" {
		HttpVarNew.Port = 80
	} else if HttpVarNew.Scheme == "https" {
		HttpVarNew.Port = 443
	}

	if parsedURL.Path != "" {
		HttpVarNew.Path = parsedURL.Path
	}

	return HttpVarNew, nil
}
