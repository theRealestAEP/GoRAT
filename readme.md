Just a quick code snippet of remote terminal access with go lang

There are two versions one that broadcasts over UDP so that the client can find the server without knowing the server IP and only knowing the client IP and no broadcasting where both the client and server IPs must be known 

to build it you can use the following commands 

broadcasting is still a WIP

Build for Windows:
```
# For 64-bit Windows
GOOS=windows GOARCH=amd64 go build -o client_windows_amd64.exe client.go
GOOS=windows GOARCH=amd64 go build -o server_windows_amd64.exe server.go

# For 32-bit Windows
GOOS=windows GOARCH=386 go build -o client_windows_386.exe client.go
GOOS=windows GOARCH=386 go build -o server_windows_386.exe server.go
```

Build for macOS:
```
# For 64-bit macOS
GOOS=darwin GOARCH=amd64 go build -o client_macos_amd64 client.go
GOOS=darwin GOARCH=amd64 go build -o server_macos_amd64 server.go

# For ARM-based macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o client_macos_arm64 client.go
GOOS=darwin GOARCH=arm64 go build -o server_macos_arm64 server.go
```

Build for Linux
```
# For 64-bit Linux
GOOS=linux GOARCH=amd64 go build -o client_linux_amd64 client.go
GOOS=linux GOARCH=amd64 go build -o server_linux_amd64 server.go

# For 32-bit Linux
GOOS=linux GOARCH=386 go build -o client_linux_386 client.go
GOOS=linux GOARCH=386 go build -o server_linux_386 server.go

# For ARM-based Linux (Raspberry Pi and others)
GOOS=linux GOARCH=arm go build -o client_linux_arm client.go
GOOS=linux GOARCH=arm go build -o server_linux_arm server.go
```
