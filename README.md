# Ping Tracker

A concurrent network scanning tool that pings IP addresses or hostnames and displays their status in real-time.

## Features

- Concurrent scanning of multiple IP addresses or hostnames
- Real-time status updates in the terminal
- Support for both ping and TCP connection checks
- Configurable timeout and retry settings

## Prerequisites

- Go 1.16 or higher
- Git

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/ping-tracker.git
cd ping-tracker
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o pingtracker
```

## Usage

### Basic Usage

To scan a range of IP addresses:
```bash
./pingtracker -range 192.168.1.1-192.168.1.255
```

To scan a list of hostnames:
```bash
./pingtracker -hosts example.com,google.com,github.com
```

### Command Line Options

- `-range`: Specify an IP range (e.g., 192.168.1.1-192.168.1.255)
- `-hosts`: Specify a comma-separated list of hostnames
- `-timeout`: Set connection timeout in seconds (default: 5)
- `-retries`: Number of retry attempts (default: 3)
- `-port`: TCP port to check (default: 80)

### Examples

Scan a specific IP range with custom timeout:
```bash
./pingtracker -range 10.0.0.1-10.0.0.255 -timeout 10
```

Scan multiple hosts with custom port:
```bash
./pingtracker -hosts example.com,google.com -port 443
```

## Testing

Run the test suite:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```
