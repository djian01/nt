# nt (Net Test)

`nt` is a Swiss Army knife for network troubleshooting and testing written in Go. It provides continuous testing of network connectivity using various protocols using subcommands, including `icmp`, `tcp`, `http`, and `dns`. Monitor round-trip times (RTT), track packet loss rates, and log high latency events with timestamps to ensure your network's reliability.

## Features

- **Subcommands for Protocols**: Use specific subcommands (`icmp`, `tcp`, `http`, `dns`) to perform tests.
- **Continuous ICMP Ping**: Monitor network latency and packet loss using ICMP echo requests.
- **TCP Connectivity Testing**: Check the availability and response times of TCP ports.
- **HTTP/HTTPS Testing**: Test web server availability and measure HTTP/HTTPS response times.
- **DNS Ping**: Verify DNS server responsiveness and resolve lookup times.
- **MTU Testing**: Determine the Maximum Transmission Unit (MTU) size to a given destination host/IP.
- **TCP Port Testing**: Test if the remote server is listening on one or multiple given TCP ports.
- **Recording and Logging**: Save test results to a CSV file for later analysis.
- **Customizable Output**: Adjust the number of rows displayed in the terminal during live tests.
- **Cross-Platform**: Compatible with Windows, macOS, and Linux.



## Installation

### Prerequisites

- Go 1.18 or higher installed on your system.

### Option 1: Install via `go install`

```bash
go install github.com/djian01/nt@latest
```

### Option 2: Build from Source

```bash
git clone https://github.com/djian01/nt.git
cd nt
go build .
```
### Option 3: Makefile (requires Docker installed in local machine)

```bash
git clone https://github.com/djian01/nt.git
cd nt
make build-linux
```

or

```bash
git clone https://github.com/djian01/nt.git
cd nt
make build-windows
```

### Option 4: Download Executable from GitHub Releases

1. Visit the [Releases](https://github.com/djian01/nt/releases) page of the repository.
2. Download the `nt_linux_amd64.tar.gz` file for Linux or the `nt_windows_amd64.zip` file for Windows


# Compile for Windows Executable
GOOS=windows GOARCH=amd64 go build -o nt.exe main.go

# Compile for macOS Executable
GOOS=darwin GOARCH=amd64 go build -o nt.exe main.go