package record

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/djian01/nt/pkg/ntPinger"
)

// Func - RecordingFunc, saving the accumulated results into CSV file
func RecordingFunc(filePath string, fileName string, bucket int, recordingChan <-chan ntPinger.Packet, wg *sync.WaitGroup) {

	// Initial the bucket
	count := 0
	writeHeader := true
	accumulatedRecords := []ntPinger.Packet{}

	// file path name
	recordingFilePathName := filepath.Join(filePath, fmt.Sprintf("[recording...]_%s", fileName))
	completedFilePathName := filepath.Join(filePath, fileName)

	// The ticker loop for CSV file write
	for {
		r, ok := <-recordingChan

		// count
		count++

		// if the recordingChan is closed, save the rest of the accumulatedRecords to CSV
		if !ok {
			// save to CSV
			err := SaveToCSV(recordingFilePathName, accumulatedRecords, writeHeader)
			if err != nil {
				panic(err)
			}
			// reset bucket
			accumulatedRecords = nil

			// change the file name, remove the "recording_" prefix
			err = os.Rename(recordingFilePathName, completedFilePathName)
			if err != nil {
				panic(err)
			}

			// clear wait group
			wg.Done()

			break

			// else adding result to accumulatedRecords
		} else {

			accumulatedRecords = append(accumulatedRecords, r)

			// if the bucket is full, Save to CSV
			if count%bucket == 0 {
				// save to CSV
				err := SaveToCSV(recordingFilePathName, accumulatedRecords, writeHeader)
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
					"RTT(ms)",
					"SendDate",
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
				// interface assertion
				pkt := recordItem.(*ntPinger.PacketICMP)
				RTT := (float64((pkt.RTT).Nanoseconds())) / 1e6

				row := []string{
					pkt.Type,                            // Ping Type
					strconv.Itoa(pkt.Seq),               // Seq
					fmt.Sprintf("%t", pkt.Status),       // Status
					pkt.DestHost,                        // DestHost
					pkt.DestAddr,                        // DestAddr
					strconv.Itoa(pkt.PayLoadSize),       // PayLoadSize
					fmt.Sprintf("%v", RTT),              // Response_Time
					pkt.SendTime.Format("2006-01-02"),   // SendDate
					pkt.SendTime.Format("15:04:05 MST"), // SendTime

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
					"RTT(ms)",
					"SendDate",
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
				// interface assertion
				pkt := recordItem.(*ntPinger.PacketTCP)
				RTT := (float64((pkt.RTT).Nanoseconds())) / 1e6

				row := []string{
					pkt.Type,                            // Ping Type
					strconv.Itoa(pkt.Seq),               // Seq
					fmt.Sprintf("%t", pkt.Status),       // Status
					pkt.DestHost,                        // DestHost
					pkt.DestAddr,                        // DestAddr
					strconv.Itoa(pkt.DestPort),          // DestPort
					strconv.Itoa(pkt.PayLoadSize),       // PayLoadSize
					fmt.Sprintf("%v", RTT),              // Response_Time
					pkt.SendTime.Format("2006-01-02"),   // SendDate
					pkt.SendTime.Format("15:04:05 MST"), // SendTime

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
		case "http":
			// Write the header if requested
			if writeHeader {
				header := []string{
					"Type",
					"Seq",
					"Status",
					"Method",
					"URL",
					"Response_Code",
					"Response_Phase",
					"Response_Time(ms)",
					"SendDate",
					"SendTime",
					"SessionSent",
					"SessionSuccess",
					"FailureRate",
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

				// interface assertion
				pkt := recordItem.(*ntPinger.PacketHTTP)
				RTT := (float64((pkt.RTT).Nanoseconds())) / 1e6

				// url
				url := ntPinger.ConstructURL(pkt.Http_scheme, pkt.DestHost, pkt.Http_path, pkt.DestPort)

				row := []string{
					pkt.Type,                             // Ping Type
					strconv.Itoa(pkt.Seq),                // Seq
					fmt.Sprintf("%t", pkt.Status),        // Status
					pkt.Http_method,                      // HTTP Method
					url,                                  // DestHost
					strconv.Itoa(pkt.Http_response_code), // Response_Code
					pkt.Http_response,                    // Response Phase
					fmt.Sprintf("%v", RTT),               // Response_Time
					pkt.SendTime.Format("2006-01-02"),    // SendDate
					pkt.SendTime.Format("15:04:05 MST"),  // SendTime

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
		case "dns":
			// Write the header if requested
			if writeHeader {
				header := []string{
					"Type",
					"Seq",
					"Status",
					"DNS_Resolver",
					"DNS_Query",
					"DNS_Response",
					"Record",
					"DNS_Protocol",
					"Response_Time(ms)",
					"SendDate",
					"SendTime",
					"PacketsSent",
					"SuccessResponse",
					"FailureRate",
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
				// interface assertion
				pkt := recordItem.(*ntPinger.PacketDNS)
				RTT := (float64((pkt.RTT).Nanoseconds())) / 1e6

				row := []string{
					pkt.Type,                            // Ping Type
					strconv.Itoa(pkt.Seq),               // Seq
					fmt.Sprintf("%t", pkt.Status),       // Status
					pkt.DestHost,                        // DNS_Resolver
					pkt.Dns_query,                       // DNS_Query
					pkt.Dns_response,                    // DNS_Response
					pkt.Dns_queryType,                   // Record
					pkt.Dns_protocol,                    // DNS_Protocol
					fmt.Sprintf("%v", RTT),              // Response_Time
					pkt.SendTime.Format("2006-01-02"),   // SendDate
					pkt.SendTime.Format("15:04:05 MST"), // SendTime
					strconv.Itoa(pkt.PacketsSent),       // PacketsSent
					strconv.Itoa(pkt.PacketsRecv),       // PacketsRecv
					fmt.Sprintf("%.2f%%", float64(pkt.PacketLoss*100)), // PacketLoss
					pkt.MinRtt.String(), // MinRtt
					pkt.AvgRtt.String(), // AvgRtt
					pkt.MaxRtt.String(), // MaxRtt
					pkt.AdditionalInfo,  // AdditionalInfo
				}

				if err := writer.Write(row); err != nil {
					return fmt.Errorf("could not write record to file: %v", err)
				}
			}
		}
	}
	return nil
}
