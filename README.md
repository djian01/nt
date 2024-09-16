# nt (Net Test)

**nt** is a versatile network testing tool designed to perform continuous and ad-hoc network tests using various protocols like ICMP, TCP, HTTP/HTTPS, and DNS. It also features an MTU scan to determine the maximum MTU for a given path. The tool can log test results to CSV files for later analysis.

## Features


# Compile for Windows Executable
GOOS=windows GOARCH=amd64 go build -o nt.exe main.go

# Compile for macOS Executable
GOOS=darwin GOARCH=amd64 go build -o nt.exe main.go