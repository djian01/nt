# Network Probe API Flows

## ICMP command

- Command: `nt icmp [flags] <host>`
- Request: host, count, payload size, timeout, interval, DF flag, display rows, and recording flag.
- Response: live `PacketICMP` rows, aggregate statistics, an optional CSV file, or a returned error.
- Frontend consumers: none; this repository provides a terminal CLI.

```mermaid
flowchart TD
    A["pkg/cmd/root/root.go<br/>RootCommand"] --> B["pkg/cmd/icmp/icmp_main.go<br/>IcmpCommandLink"]
    B --> C["pkg/cmd/icmp/icmp_main.go<br/>IcmpCommandMain"]
    C --> D["pkg/ntPinger/Functions.go<br/>NewPinger"]
    D --> E["pkg/ntPinger/Struct_Method.go<br/>Pinger.Resolve"]
    C --> F["pkg/ntPinger/Struct_Method.go<br/>Pinger.Run"]
    F --> G["pkg/ntPinger/ProbingICMP_Run.go<br/>icmpProbingRun"]
    G --> H["pkg/ntPinger/ProbingRun.go<br/>runProbes"]
    H --> I["pkg/ntPinger/ProbingICMP.go<br/>IcmpProbing"]
    I -->|"macOS/Linux standard probe"| J["pkg/ntPinger/ProbingICMP.go<br/>icmpNativeProbing"]
    I -->|"DF, Windows, or native unavailable"| K["pkg/ntPinger/ProbingICMP_{OS}.go<br/>icmpCommandProbing"]
    H --> L["pkg/terminalOutput/output_main.go<br/>OutputFunc"]
    H --> M["pkg/record/record_main.go<br/>RecordingFunc"]
```

`IcmpProbing` keeps its existing public signature. Standard probes prefer the
unprivileged native transport on macOS and Linux. A socket availability or
permission failure selects the command compatibility path. DF and MTU probes
select that path directly because their socket options are platform-specific.

## Shared finite-probe lifecycle

- Commands: `nt icmp`, `nt tcp`, `nt http`, and `nt dns`.
- Request validation: nonnegative count and interval, positive timeout and display rows, valid payload and port ranges, and `udp` or `tcp` for DNS.
- Response lifecycle: each successful probe is emitted once; fatal errors stop the producer; output and optional recording finish before the command returns.
- Frontend consumers: none.

```mermaid
flowchart TD
    A["pkg/cmd/{icmp,tcp,http,dns}/*_main.go<br/>CommandMain"] --> B["pkg/ntPinger/Functions.go<br/>NewPinger and validateInputVars"]
    B --> C["pkg/ntPinger/Struct_Method.go<br/>Pinger.Run"]
    C --> D["pkg/ntPinger/ProbingRun.go<br/>runProbes"]
    D --> E["protocol probe function"]
    E --> D
    D --> F["pkg/terminalOutput/output_main.go<br/>OutputFunc"]
    D --> G["pkg/record/record_main.go<br/>RecordingFunc and SaveToCSV"]
    F --> H["CommandMain waits for output completion"]
    G --> H
```

## HTTP command

- Command: `nt http [flags] <http-or-https-url>`
- Request: URL, method, allowed status codes, count, timeout, interval, optional proxy, display rows, and recording flag.
- Response: live `PacketHTTP` rows, aggregate statistics, an optional CSV file, or a returned error.
- Frontend consumers: none.

```mermaid
flowchart TD
    A["pkg/cmd/http/http_main.go<br/>HttpCommandLink"] --> B["pkg/cmd/http/http_main.go<br/>ParseURL and ParseStatusCodes"]
    B --> C["pkg/cmd/http/http_main.go<br/>HttpCommandMain"]
    C --> D["pkg/ntPinger/Functions.go<br/>NewPinger"]
    D --> E["pkg/ntPinger/Struct_Method.go<br/>Pinger.Run"]
    E --> F["pkg/ntPinger/ProbingHTTP.go<br/>httpProbingRun"]
    F --> G["pkg/ntPinger/ProbingRun.go<br/>runProbes"]
    G --> H["pkg/ntPinger/ProbingHTTP.go<br/>HttpProbing"]
```

`ParseURL` accepts only HTTP and HTTPS URLs with a hostname and a port from 1
through 65535. It preserves the escaped request path and query string.

## TCP scan command

- Command: `nt tcpscan [flags] <host> <port-or-range> ...`
- Request: host, one to fifty valid TCP ports, timeout, and recording flag.
- Response: one immutable result per port, terminal scan statistics, an optional CSV file, or a returned error.
- Frontend consumers: none.

```mermaid
flowchart TD
    A["pkg/cmd/tcpscan/tcpscan_main.go<br/>TcpScanCommandLink"] --> B["pkg/cmd/tcpscan/tcpscan_main.go<br/>TcpScanCommandMain"]
    B --> C["pkg/ntPinger/Functions.go<br/>ResolveDestHost"]
    B --> D["pkg/ntScaner/scanTCP.go<br/>ScanTcpWorkerResults"]
    D --> E["pkg/ntScaner/scanTCP.go<br/>scanTCPPort"]
    E -->|"three concurrent attempts"| F["pkg/ntPinger/ProbingTCP.go<br/>TcpProbing"]
    D --> G["pkg/cmd/tcpscan/tcpscan_main.go<br/>single-owner result collection"]
    G --> H["pkg/terminalOutput/output_tcpScanTablePrint.go<br/>ScanTablePrint"]
    G --> I["pkg/record/record_tcpScanTablePrint.go<br/>TcpScan_Recording"]
```

The command owns the shared display table. Workers receive values and return
values, preventing concurrent reads and writes to the table. The original
`ScanTcpWorker` function remains available for compatibility.

## MTU command

- Command: `nt mtu [flags] <host>`
- Request: host and a ceiling from 29 through 65535 bytes.
- Response: the largest confirmed IPv4 MTU at or below the ceiling, or an error when the minimum probe receives no reply.
- Frontend consumers: none.

```mermaid
flowchart TD
    A["pkg/cmd/mtu/mtu.go<br/>MtuCommandLink"] --> B["pkg/cmd/mtu/mtu.go<br/>MtuCommandMain"]
    B --> C["pkg/ntPinger/Functions.go<br/>ResolveDestHost"]
    B --> D["pkg/ntScaner/scanMTU.go<br/>ScanMTURun"]
    D --> E["pkg/ntScaner/scanMTU.go<br/>findLargestPayload"]
    E --> F["pkg/ntPinger/ProbingICMP.go<br/>IcmpProbing with DF"]
    F --> G["pkg/ntPinger/ProbingICMP_{OS}.go<br/>icmpCommandProbing"]
```

The search first verifies reachability with a zero-byte payload and then uses a
bounded binary search. The payload bounds are `0` through `ceiling - 28`, where
28 bytes represents the IPv4 and ICMP headers.
