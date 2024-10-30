# nt (Net Test) - A versatile network testing tool that supports ICMP, HTTP(s), DNS, TCP, and MTU testing

`nt` is a versatile tool for network troubleshooting and testing written in Go. It provides continuous testing of network connectivity using various protocols using subcommands, including `icmp`, `tcp`, `http`, and `dns`. Monitor round-trip times (RTT), track packet loss rates, and log high latency events with timestamps to ensure your network's reliability.

## Table of Contents
- [Features](#features)
- [Installation](#installation)
  - [Option 1 - go install](#option-1-install-via-go-install)
  - [Option 2 - build from source](#option-2-build-from-source)
  - [Option 3 - Makefile](#option-3-makefile-docker--make-are-required-on-your-system)
  - [Option 4 - Download from Release](#option-4-download-executable-from-github-releases)
- [Usage](#usage)
  - [Main Command](#main-command)
  - [ICMP Ping Sub-Command](#icmp-sub-command-required-privilege-mode-in-linux)
  - [TCP Ping Sub-Command](#tcp-sub-command)
  - [HTTP Ping Sub-Command](#http-sub-command)
  - [DNS Ping Sub-Command](#dns-sub-command)
  - [MTU Discovery Sub-Command](#mtu-sub-command)
  - [TCP Scan Sub-Command](#tcp-scan-sub-command)
- [Logging and Recording](#logging-and-recording)
- [License](#license)  
- [Contact](#contact)  

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
### Option 3: Makefile (Docker & Make are required on your system)

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
nt [flags] <sub-command: icmp/tcp/http/dns/mtu/tcpscan> [args]

```

#### Global Options
- `-r`:   **Enable Recording**
 Save the test results to a CSV file for future analysis.

- `-p <number>`:   **Rows Displayed**
  Specify the number of rows displayed in the terminal during live tests. Default is `10`.



### ICMP Sub-Command (required privilege mode in Linux)

#### ICMP Options
- `-c`:   **ICMP Ping Count**  
  Number of ICMP ping requests to send. Default is `0`, which means it will run non-stop until interrupted with `Ctrl+C`.

- `-d`:    *ICMP Ping Do Not Fragment**  
  Sets the "Do Not Fragment" flag in the ICMP request. Default is `false`.

- `-h`:   **Help**  
  Display help information for the `icmp` subcommand.

- `-i`:   **ICMP Ping Interval**  
  Interval between ICMP ping requests in seconds. Default is `1` second.

- `-s`:   **ICMP Ping Payload Size**  
  Size of the ICMP ping payload in bytes. Default is `32` bytes.

- `-t`:   **ICMP Ping Timeout**  
  Timeout for each ICMP ping request in seconds. Default is `4` seconds.

#### Example 1: ICMP continuous ping to "google.com" with recording enabled

```bash
nt -r icmp google.com

```

#### Example 2: ICMP ping to "10.2.3.10" with count: 10, interval: 2 sec,  payload 48 bytes

```bash
nt icmp -c 10 -i 2 -s 48 10.2.3.10

```


### TCP Sub-Command

#### TCP Options
- `-c`:    **TCP Ping Count**  
  Number of TCP ping requests to send. Default is `0`, which means it will run non-stop until interrupted with `Ctrl+C`.

- `-h`:   **Help**  
  Display help information for the `tcp` subcommand.

- `-i`:   **TCP Ping Interval**  
  Interval between TCP ping requests in seconds. Default is `1` second.

- `-s`:   **TCP Ping Payload Size**  
  Size of the TCP ping payload in bytes. Default is `0` bytes (no payload).

- `-t`:   **TCP Ping Timeout**  
  Timeout for each TCP ping request in seconds. Default is `4` seconds.


#### Example 1: TCP ping to "google.com:443" with recording enabled

```bash
nt -r tcp google.com 443

```

#### Example 2: TCP ping to "10.2.3.10:22" with count: 10 and interval: 2 sec

```bash
nt tcp -c 10 -i 2 10.2.3.10 22

```

### HTTP Sub-Command
**Note:**  
- If custom port is not specificed, the test will use the default ports (HTTP 80, HTTPS 443).
- The default interval for HTTP Ping is 5 seconds.

#### HTTP Options
- `-c`:   **HTTP Ping Count**  
  Number of HTTP ping requests to send. Default is `0`, which means it will run non-stop until interrupted with `Ctrl+C`.

- `-h`:   **Help**  
  Display help information for the `http` subcommand.

- `-i`:   **HTTP Ping Interval**  
  Interval between HTTP ping requests in seconds. Default is `5` seconds.

- `-m`:   **HTTP Ping Method**  
  HTTP request method to use (e.g., `GET`, `POST`). Default is `"GET"`.

- `-t`:   **HTTP Ping Timeout**  
  Timeout for each HTTP ping request in seconds. Default is `4` seconds.


#### Example 1: HTTP ping to "https://google.com" with recording enabled (With default values: Port-443, Method-GET, Count-0, Interval-5s, Timeout-4s)

```bash
nt -r http https://google.com

```

#### Example 2: HTTP ping to POST "http://10.2.3.10:8080/token" with count: 10 and interval: 2 sec

```bash
nt http -c 10 -i 2 -m POST http://10.2.3.10:8080/token

```

### DNS Sub-Command

#### DNS Options
- `-c`:   **DNS Ping Count**  
  Number of DNS ping requests to send. Default is `0`, which means it will run non-stop until interrupted with `Ctrl+C`.

- `-h`:   **Help**  
  Display help information for the `dns` subcommand.

- `-i`:   **DNS Ping Interval**  
  Interval between DNS ping requests in seconds. Default is `1` second.

- `-o`:   **DNS Ping Protocol Type**  
  Protocol to use for DNS queries (e.g., `udp`, `tcp`). Default is `"udp"`.

- `-t`:   **DNS Ping Timeout**  
  Timeout for each DNS ping request in seconds. Default is `4` seconds.


#### Example 1: DNS ping to "8.8.8.8" with query "google.com" and have recording enabled

```bash
nt -r dns 8.8.8.8 google.com

```

#### Example 2: DNS ping to "4.2.2.2" with query "abc.com" with count: 10 and interval: 2 sec

```bash
nt dns -c 10 -i 2 4.2.2.2 abc.com

```

### MTU Sub-Command

#### MTU Options
- `-s`:   **Ceiling Test Size**  
  The maximum MTU size to test, in bytes. Default is `1500` bytes.

- `-h`:   **Help**  
  Display help information for the `mtu` subcommand.

#### Example 1: MTU check for destination google.com
```bash
nt mtu google.com

```

#### Example 2: MTU check for destination 192.168.1.10 with user defined ceiling test size 9000 set (for Jumbo Frame enabled environment)
```bash
nt mtu -s 9000 192.168.1.10

```

### TCP SCAN Sub-Command
**Note:**  
- The maximum number of tested ports for each command run is 50.

#### TCP SCAN  Options
- `-t`:   **TCP Ping test Timeout **  
  TCP Ping test Timeout (default: 4 sec) (default 4)

- `-h`:   **Help**  
  Display help information for the `mtu` subcommand.

#### Example 1: TCP Scan to "10.123.1.10" for port "80, 443, 8080 & 1500-1505" with recording enabled
```bash
nt -r tcpscan 10.123.1.10 80 443 8080 1500-1505

```

#### Example 2: TCP SCAN to "10.2.3.10" for port "22, 1522-1525 & 8433" with custom timeout: 5 sec
```bash
nt tcpscan -t 5 10.2.3.10 22 1522-1525 8433

```


## Logging and Recording
When the `-r` option is enabled, all test results are saved to a CSV file in the same directory as the executable. The CSV files are named using the format:

```bash
Record_<test type>_<test target host>_<timestamp>.csv

```

- **Example**: `Record_icmp_google.com_20211012T101530.csv` 

The CSV file captures the detailed test results based on the test type.



## License

This project is licensed under the [MIT License](https://opensource.org/licenses/MIT).


## Contact

- **Email**: [dennis.jian@packetstreams.net](mailto:dennis.jian@packetstreams.net)
- **GitHub**: [djian01](https://github.com/djian01)