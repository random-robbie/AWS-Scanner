# AWS-Scanner

A high-performance security tool for discovering AWS S3 buckets and CloudFront distributions on target websites.

## Overview

AWS-Scanner efficiently scans a list of URLs to identify exposed AWS resources including:
- Amazon S3 buckets
- CloudFront distributions

The scanner outputs clean CSV files with discovered resources for further analysis.

## Features

- ✅ **Auto-detection of URL protocols** (http/https)
- ✅ **Smart JavaScript filtering** to avoid parsing messy inline scripts
- ✅ **Response size limiting** (5MB) to prevent memory issues
- ✅ **Unified regex engine** for efficient pattern matching
- ✅ **Deduplication** of discovered resources
- ✅ **Clean CSV output** with normalized HTTPS URLs
- ✅ **Concurrent scanning** with configurable timeouts

## Installation

### Requirements
- Go 1.16 or later

### Building from Source

```bash
# Clone the repository
git clone https://github.com/random-robbie/AWS-Scanner.git
cd AWS-Scanner

# Download dependencies
go mod tidy

# Build the binary
go build -o aws-scanner main.go
```

### Quick Run (Without Building)

```bash
go run main.go --list list.txt
```

### Using Docker

Build the Docker image:

```bash
docker build -t aws-scanner .
```

Run the scanner with Docker:

```bash
# Mount your URL list and output directory
docker run -v $(pwd)/list.txt:/app/list.txt -v $(pwd)/output:/app/output aws-scanner --list list.txt

# The CSV files will be saved to ./output directory
```

Or use docker-compose (create a `docker-compose.yml`):

```yaml
version: '3.8'
services:
  aws-scanner:
    build: .
    volumes:
      - ./list.txt:/app/list.txt
      - ./output:/app/output
    command: --list list.txt
```

Then run:

```bash
docker-compose up
```

## Usage

### Basic Scan

```bash
./aws-scanner --list list.txt
```

### Input Format

Create a text file (`list.txt`) with one URL per line:

```
github.com
https://www.example.com
http://test.com
```

**Note:** URLs can be provided with or without protocol prefixes (http:// or https://). The scanner will automatically detect and normalize them.

### Output

The scanner generates two CSV files:

- **`s3-bucket.csv`** - Discovered S3 buckets
- **`cloudfront.csv`** - Discovered CloudFront distributions

**Format:** `source_url,discovered_resource`

Example:
```csv
https://github.com/,https://github-cloud.s3.amazonaws.com
https://example.com/,https://example-assets.cloudfront.net
```

## Recent Improvements

- [x] Unified regex function for cleaner codebase
- [x] Auto-detection of URL protocols
- [x] JavaScript and script tag filtering
- [x] Response body size limiting
- [x] Normalized HTTPS-only output
- [x] Duplicate resource removal

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues for bugs and feature requests.

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

See [LICENSE](LICENSE) file for details.

## Credits

Special thanks to [Glove](https://github.com/Glove) for contributions.

## Hosting

Use a VPS from DigitalOcean:

[![DigitalOcean Referral Badge](https://web-platforms.sfo2.cdn.digitaloceanspaces.com/WWW/Badge%201.svg)](https://www.digitalocean.com/?refcode=e22bbff5f6f1&utm_campaign=Referral_Invite&utm_medium=Referral_Program&utm_source=badge)

---

**Disclaimer:** This tool is intended for security research and authorized testing only. Always obtain proper authorization before scanning any websites or networks you do not own.
