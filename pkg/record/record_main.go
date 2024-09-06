package record

import (
	"encoding/csv"
	"fmt"
	"nt/pkg/ntPinger"
	"os"
	"strconv"
	"sync"
)

// Func - RecordingFunc, saving the accumulated results into CSV file
func RecordingFunc(filePath string, bucket int, recordingChan <-chan ntPinger.Packet, wg *sync.WaitGroup) {

	// Initial the bucket
	count := 0
	writeHeader := true
	accumulatedRecords := []ntPinger.Packet{}

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

			// clear wait group
			wg.Done()

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

func SaveToCSV(filePath string, accumulatedRecords []ntPinger.Packet, writeHeader bool) error {

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
		// else Save to CSV based on Type
	} else {
		switch accumulatedRecords[0].GetType() {
		case "icmp":
			// Write the header if requested
			if writeHeader {
				header := []string{
					"Type",
					"Seq",
					"Status",
					"DestHost",
					"DestAddr",
					"PayLoadSize",
					"RTT",
					"SendTime",
					"PacketsSent",
					"PacketsRecv",
					"PacketLoss",
					"MinRtt",
					"AvgRtt",
					"MaxRtt",
					"AdditionalInfo",
				}

				err := writer.Write(header)

				if err != nil {
					return fmt.Errorf("could not write header to file: %v", err)
				}
			}

			// Write each struct to the file
			for _, recordItem := range accumulatedRecords {
				pkt := recordItem.(*ntPinger.PacketICMP)
				row := []string{
					pkt.Type,                      // Ping Type
					strconv.Itoa(pkt.Seq),         // Seq
					fmt.Sprintf("%t", pkt.Status), // Status
					pkt.DestHost,                  // DestHost
					pkt.DestAddr,                  // DestAddr
					strconv.Itoa(pkt.PayLoadSize), // PayLoadSize
					(pkt.RTT).String(),            // RTT
					pkt.SendTime.Format("2006-01-02 15:04:05"), // SendTime

					strconv.Itoa(pkt.PacketsSent),                      // PacketsSent
					strconv.Itoa(pkt.PacketsRecv),                      // PacketsRecv
					fmt.Sprintf("%.2f%%", float64(pkt.PacketLoss*100)), // PacketLoss
					pkt.MinRtt.String(),                                // MinRtt
					pkt.AvgRtt.String(),                                // AvgRtt
					pkt.MaxRtt.String(),                                // MaxRtt
					pkt.AdditionalInfo,                                 // AdditionalInfo
				}

				if err := writer.Write(row); err != nil {
					return fmt.Errorf("could not write record to file: %v", err)
				}
			}
		case "tcp":
			// Write the header if requested
			if writeHeader {
				header := []string{
					"Type",
					"Seq",
					"Status",
					"DestHost",
					"DestAddr",
					"DestPort",
					"PayLoadSize",
					"RTT",
					"SendTime",
					"PacketsSent",
					"PacketsRecv",
					"PacketLoss",
					"MinRtt",
					"AvgRtt",
					"MaxRtt",
					"AdditionalInfo",
				}

				err := writer.Write(header)

				if err != nil {
					return fmt.Errorf("could not write header to file: %v", err)
				}
			}

			// Write each struct to the file
			for _, recordItem := range accumulatedRecords {
				pkt := recordItem.(*ntPinger.PacketTCP)
				row := []string{
					pkt.Type,                                   // Ping Type
					strconv.Itoa(pkt.Seq),                      // Seq
					fmt.Sprintf("%t", pkt.Status),              // Status
					pkt.DestHost,                               // DestHost
					pkt.DestAddr,                               // DestAddr
					strconv.Itoa(pkt.DestPort),                 // DestPort
					strconv.Itoa(pkt.PayLoadSize),              // PayLoadSize
					(pkt.RTT).String(),                         // RTT
					pkt.SendTime.Format("2006-01-02 15:04:05"), // SendTime

					strconv.Itoa(pkt.PacketsSent),                      // PacketsSent
					strconv.Itoa(pkt.PacketsRecv),                      // PacketsRecv
					fmt.Sprintf("%.2f%%", float64(pkt.PacketLoss*100)), // PacketLoss
					pkt.MinRtt.String(),                                // MinRtt
					pkt.AvgRtt.String(),                                // AvgRtt
					pkt.MaxRtt.String(),                                // MaxRtt
					pkt.AdditionalInfo,                                 // AdditionalInfo
				}

				if err := writer.Write(row); err != nil {
					return fmt.Errorf("could not write record to file: %v", err)
				}
			}
		}
	}
	return nil
}
