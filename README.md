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
- **Cross-Platform**: Compatible with Windows and Linux.



## Installation

### Prerequisites

- For Option 1 & 2, Go 1.18 or higher is required on your system.

### Option 1: Install via `go install`

The compiled executable will be placed in `$GOPATH/bin`

```bash
go install github.com/djian01/nt@latest
```

### Option 2: Build from Source

```bash
git clone https://github.com/djian01/nt.git
cd nt
go build .
```
### Option 3: Makefile (requires Docker & Make are required on your system)

The compiled executable will be placed in `\executable\` inside the code folder

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
2. Download the `nt_linux_amd64_x.x.x.tar.gz` file for Linux or the `nt_windows_amd64_x.x.x.zip` file for Windows

## Usage

### Main Command

```bash
nt [flags] <sub-command: icmp/tcp/http/dns/mtu> [args]

```

#### Global Options
- `-r`: **Enable Recording**
 Save the test results to a CSV file for future analysis.

- `-p <number>`: **Rows Displayed**
  Specify the number of rows displayed in the terminal during live tests. Default is `10`.