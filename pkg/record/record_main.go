package record

import (
	"encoding/csv"
	"fmt"
	"nt/pkg/sharedStruct"
	"os"
	"reflect"
	"strconv"
)

// Func - RecordingFunc, saving the accumulated results into CSV file
func RecordingFunc(Type string, filePath string, bucket int, recordingChan <-chan sharedStruct.NtResult) {

	// Initial the bucket
	count := 0
	writeHeader := true
	accumulatedRecords := []sharedStruct.NtResult{}

	// The ticker loop for CSV file write
	for {
		r, ok := <-recordingChan

		// count
		count++

		// if the recordingChan is closed, save the rest of the accumulatedRecords to CSV
		if !ok {
			// save to CSV
			err := SaveToCSV(filePath, accumulatedRecords, writeHeader)
			if err != nil {
				panic(err)
			}
			// reset bucket
			accumulatedRecords = nil

			break

			// else adding result to accumulatedRecords
		} else {

			accumulatedRecords = append(accumulatedRecords, r)

			// if the bucket is full, Save to CSV
			if count%bucket == 0 {
				// save to CSV
				err := SaveToCSV(filePath, accumulatedRecords, writeHeader)
				if err != nil {
					panic(err)
				}

				// set header Flag
				writeHeader = false

				// reset bucket
				accumulatedRecords = nil
			}

		}
	}

}

func SaveToCSV(filePath string, accumulatedRecords []sharedStruct.NtResult, writeHeader bool) error {

	// Open or create the file with append mode and write-only access
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// if accumulatedRecords is empty
	if len(accumulatedRecords) == 0 {
		return nil
	} else {
		switch accumulatedRecords[0].Type {
		case "icmp":
			// Write the header if requested
			if writeHeader {
				header := []string{
					"Seq",
					"Status",
					"HostName",
					"IP",
					"Size",
					"RTT",
					"Timestamp",
					"Type",
					"PacketsSent",
					"PacketsRecv",
					"PacketLoss",
					"MinRtt",
					"AvgRtt",
					"MaxRtt"}

				err := writer.Write(header)

				if err != nil {
					return fmt.Errorf("could not write header to file: %v", err)
				}
			}

			// Write each struct to the file
			for _, recordItem := range accumulatedRecords {
				row := []string{
					strconv.Itoa(recordItem.Seq),         // Seq
					recordItem.Status,                    // Status
					recordItem.HostName,                  // Hostname
					recordItem.IP,                        // IP
					strconv.Itoa(recordItem.Size),        // Size
					(recordItem.RTT).String(),            // RTT
					recordItem.Timestamp,                 // TimeStamp
					recordItem.Type,                      // Ping Type
					strconv.Itoa(recordItem.PacketsSent), // PacketsSent
					strconv.Itoa(recordItem.PacketsRecv), // PacketsRecv
					fmt.Sprintf("%.2f%%", float64(recordItem.PacketLoss*100)), // PacketLoss
					recordItem.MinRtt.String(),                                // MinRtt
					recordItem.AvgRtt.String(),                                // AvgRtt
					recordItem.MaxRtt.String(),                                // MaxRtt
				}

				if err := writer.Write(row); err != nil {
					return fmt.Errorf("could not write record to file: %v", err)
				}
			}
		}
	}
	return nil
}

// GetFieldValueByIndex returns the value of a struct field at the specified index.
func GetFieldValueByIndex(obj interface{}, index int) (interface{}, error) {
	v := reflect.ValueOf(obj)

	// Ensure obj is a struct
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %s", v.Kind())
	}

	// Ensure the index is within the valid range
	if index < 0 || index >= v.NumField() {
		return nil, fmt.Errorf("index %d out of range", index)
	}

	// Get the field value at the specified index
	field := v.Field(index)

	return field.Interface(), nil
}
