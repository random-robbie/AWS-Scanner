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
- ✅ **Request timeouts** (30 seconds per site) to avoid hangs

## Installation

### Requirements
- Go 1.24 or later

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
# Mount your URL list and output directory, and set the working directory
# to the mounted output folder so the CSV files are written there
docker run -v $(pwd)/list.txt:/data/list.txt -v $(pwd)/output:/data -w /data aws-scanner --list list.txt

# The CSV files (s3-bucket.csv and cloudfront.csv) will be saved to ./output
```

Or use docker-compose (create a `docker-compose.yml`):

```yaml
version: '3.8'
services:
  aws-scanner:
    build: .
    working_dir: /data
    volumes:
      - ./list.txt:/data/list.txt
      - ./output:/data
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

## Author

Robbie Wiggins — [What-Security](https://www.what-security.co.uk) · [GitHub](https://github.com/random-robbie)

**Disclaimer:** This tool is intended for security research and authorised testing only. Always obtain proper authorisation before scanning any websites or networks you do not own.
