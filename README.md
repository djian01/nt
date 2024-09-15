# nt
Net Test Tool


# Compile for Windows Executable
GOOS=windows GOARCH=amd64 go build -o nt.exe main.go

# Compile for macOS Executable
GOOS=darwin GOARCH=amd64 go build -o nt.exe main.go