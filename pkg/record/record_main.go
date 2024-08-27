package record

import (
	"encoding/csv"
	"fmt"
	"nt/pkg/sharedStruct"
	"os"
	"time"
)

// Func - reportFunc, saving the accumulated results into CSV file
func RecordFunc(Type string, filePath string, accumulatedRecords *[]sharedStruct.NtResult, ticker <-chan time.Time) error {

	// open / create CSV file
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// The ticker loop for CSV file write
	for {
		_, ok := <-ticker

		// if ticker channel is closed break the loop
		if !ok {
			break
		}

		// if ticker channel is open and ticker ticks
		for _, record := range *accumulatedRecords {
			// ICMP recording
			if Type == "icmp" {
				row := []string{
					fmt.Sprintf("%d", record.ID),
					record.Name,
					fmt.Sprintf("%.2f", record.Value),
				}
				if err := writer.Write(row); err != nil {
					return fmt.Errorf("failed to write record to csv: %v", err)
				}
			}
		}
		// Clear the slice (accumulatedRecords) after writing to CSV file
		*accumulatedRecords = nil
	}

	return nil
}
