package record

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/djian01/nt/pkg/ntScaner"
	"github.com/fatih/color"
)

func TcpScan_Recording(filePath string, PortsTable []ntScaner.TcpScanPort) error {

	// Open or create the file with append mode and write-only access
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write header
	header := []string{
		"DestHost",
		"DestAddr",
		"Test_Port",
		"Status",
	}

	err = writer.Write(header)

	if err != nil {
		return fmt.Errorf("could not write header to file: %v", err)
	}

	// write CSV Body
	for _, Port := range PortsTable {

		// if Port is NOT empty
		if Port.DestAddr != "" {
			// check status
			Status := ""
			if Port.Status == 2 {
				Status = "TRUE"
			} else {
				Status = "FALSE"
			}

			row := []string{
				Port.DestHost,           // Destination Host
				Port.DestAddr,           // Destination Address
				strconv.Itoa(Port.Port), // Tested Port
				Status,                  // Test Result
			}

			if err := writer.Write(row); err != nil {
				return fmt.Errorf("could not write record to file: %v", err)
			}
		}
	}

	fmt.Printf("Recording CSV file is saved at: %s\n", color.GreenString(filePath))

	return nil
}
