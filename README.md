 # tracentp

A Network Time Protocol (NTP) tracing tool written in Go. Similar to `traceroute` but specifically designed for NTP servers, it traces the path to NTP stratum servers and displays detailed timing information.

## Features

- **NTP Tracing**: Trace the path to NTP stratum servers
- **Detailed Timing**: Display clock offset, root distance, and RTT information
- **Configurable**: Customizable timeout, port, and count parameters
- **Stratum Detection**: Automatically stops when reaching stratum 1 servers
- **Template Output**: Customizable output format using Go templates

## Installation

### Prerequisites

- Go 1.23 or later

### Build from source

```bash
git clone https://github.com/mengzhuo/tracentp.git
cd tracentp
go build -o tracentp main.go
```

### Install globally

```bash
go install github.com/mengzhuo/tracentp@latest
```

## Usage

### Basic Usage

```bash
# Trace to a specific NTP server
./tracentp pool.ntp.org

# Trace with custom port
./tracentp -p 123 time.google.com

# Trace with custom timeout
./tracentp -t 5s time.nist.gov

# Limit the number of hops
./tracentp -c 5 time.windows.com
```

### Command Line Options

| Option | Description | Default |
|--------|-------------|---------|
| `-t, --timeout` | Timeout duration for each query | `3s` |
| `-p, --port` | NTP server port | `123` |
| `-c, --count` | Stop after COUNT replies | `16` |
| `--version` | Show version information | - |

### Output Format

The tool outputs detailed information for each hop in the NTP trace:

```
OK from pool.ntp.org:seq=1 stratum=2 offset=0.000123 distance=0.001234 RTT=0.005678 ref=time.google.com
OK from time.google.com:seq=2 stratum=1 offset=0.000456 distance=0.000789 RTT=0.003456 ref=GPS
```

### Output Fields

- **Validate**: Status of the NTP response
- **Address**: NTP server address
- **Seq**: Sequence number of the hop
- **Stratum**: NTP stratum level (1 = primary reference)
- **ClockOffset**: Clock offset in seconds
- **RootDistance**: Root distance in seconds
- **RTT**: Round-trip time in seconds
- **ReferenceString**: Reference identifier

## Examples

### Trace to pool.ntp.org

```bash
./tracentp pool.ntp.org
```

### Trace with custom parameters

```bash
./tracentp -t 10s -p 123 -c 10 time.google.com
```

### Using as a ping-like tool

```bash
./tracentp -c 1 time.nist.gov
```

## How It Works

1. **Initial Query**: Sends an NTP query to the specified server
2. **Response Analysis**: Extracts timing information and reference server
3. **Recursive Tracing**: Uses the reference server as the next target
4. **Stratum Detection**: Stops when reaching a stratum 1 server
5. **Output Formatting**: Displays results using a customizable template

## Dependencies

- [github.com/alexflint/go-arg](https://github.com/alexflint/go-arg) - Command line argument parsing
- [github.com/beevik/ntp](https://github.com/beevik/ntp) - NTP client implementation

## Development

### Prerequisites

- Go 1.23 or later
- Make (optional, for using Makefile targets)

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
make test-coverage

# Run benchmarks
make bench
```

### Building

```bash
# Build for current platform
make build

# Build for multiple platforms
make build-all

# Install globally
make install
```

### Code Quality

```bash
# Format code
make format

# Lint code
make lint

# Clean build artifacts
make clean
```

### Cross-compilation

```bash
# For Linux
GOOS=linux GOARCH=amd64 go build -o tracentp-linux-amd64 main.go

# For Windows
GOOS=windows GOARCH=amd64 go build -o tracentp-windows-amd64.exe main.go

# For macOS
GOOS=darwin GOARCH=amd64 go build -o tracentp-darwin-amd64 main.go
```

### Release Process

The project uses GoReleaser for automated releases. To create a release:

1. **Create a new tag:**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. **Test release locally:**
   ```bash
   make release-dry-run
   ```

3. **GitHub Actions will automatically:**
   - Build binaries for multiple platforms
   - Create a GitHub release
   - Upload artifacts
   - Generate checksums

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Acknowledgments

- Inspired by traditional traceroute tools
- Built with the excellent NTP library by [beevik](https://github.com/beevik/ntp)
- Command line parsing powered by [go-arg](https://github.com/alexflint/go-arg)