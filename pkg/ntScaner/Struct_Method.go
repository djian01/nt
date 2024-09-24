package ntScaner

type TcpScanPort struct {
	ID       int
	Port     int
	Status   int // 0: Not_Requited (not visible), 1: Not_Check (Grey), 2: Success (Green), 3: Failed (Red)
	Timeout  int
	DestHost string
	DestAddr string
	// Status update
	AdditionalInfo string
}
