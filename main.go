package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/fatih/color"
)

var (
	list     = flag.String("list", "list.txt", "file of URLs to scan, one per line (Default: list.txt)")
	concurr  = flag.Int("c", 10, "number of concurrent workers")
	timeout  = flag.Duration("timeout", 30*time.Second, "per-request timeout")
	insecure = flag.Bool("insecure", true, "skip TLS certificate verification")

	s3path = "s3-bucket.csv"
	cfpath = "cloudfront.csv"

	s3pathMu sync.Mutex
	cfpathMu sync.Mutex

	awsurls = []string{
		"http://s3.amazonaws.com/([-A-z0-9.]+)",
		"https://s3.amazonaws.com/([-A-z0-9.]+)",
		"//s3.amazonaws.com/([-A-z0-9.]+)",
		"http://([-A-z0-9.]+).amazonaws.com",
		"https://([-A-z0-9.]+).amazonaws.com",
		"//([-A-z0-9.]+).amazonaws.com",
		"http://(s3.|s3-)[a-zA-Z0-9-]*.amazonaws.com/([-A-z0-9.]+)",
		"https://(s3.|s3-)[a-zA-Z0-9-]*.amazonaws.com/([-A-z0-9.]+)",
		"//(s3.|s3-)[a-zA-Z0-9-]*.amazonaws.com/([-A-z0-9.]+)",
	}

	cloudfronturls = []string{
		"https://(.+?).cloudfront.net/",
		"http://(.+?).cloudfront.net/",
		"//(.+?).cloudfront.net/",
	}

	// Precompiled once; the patterns are fixed at startup.
	awsRegexes        = compileAll(awsurls)
	cloudfrontRegexes = compileAll(cloudfronturls)

	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/156.0.0.0 Safari/537.36"
)

func compileAll(patterns []string) []*regexp.Regexp {
	out := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		out = append(out, regexp.MustCompile(p))
	}
	return out
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return err != nil
}

// readLines reads a whole file into memory and stores it as an array of lines.
func readLines(path string) (lines []string, err error) {
	var (
		file   *os.File
		part   []byte
		prefix bool
	)
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

// appendCSV appends a "url,finding" line to path. Each output file has its own
// mutex so concurrent workers cannot interleave writes.
func appendCSV(path string, url, finding string, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()

	fileHandle, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if isError(err) {
		return
	}
	defer fileHandle.Close()

	writer := bufio.NewWriter(fileHandle)
	fmt.Fprintln(writer, fmt.Sprintf("%s,%s", url, finding))
	writer.Flush()
}

func normalizeURL(input string) string {
	// Check if URL already has a protocol
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		return input
	}
	// Default to https if no protocol specified
	return fmt.Sprintf("https://%s", input)
}

func cleanResponse(body string) string {
	// Remove script tags and their content to avoid messy JS
	scriptRe := regexp.MustCompile(`(?is)<script[^>]*>.*?</script>`)
	cleaned := scriptRe.ReplaceAllString(body, "")

	// Remove inline script handlers
	inlineScriptRe := regexp.MustCompile(`(?i)\s+on\w+\s*=\s*["'][^"']*["']`)
	cleaned = inlineScriptRe.ReplaceAllString(cleaned, "")

	return cleaned
}

// scanURL fetches one URL and reports any S3/CloudFront references found.
// Returns true when at least one finding was written.
func scanURL(client *http.Client, rawURL string, found *atomic.Int64) {
	target := normalizeURL(strings.TrimSpace(rawURL))
	if target == "" {
		return
	}

	color.HiWhite("[*] Scanning Website %s\n", target)

	req, err := http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		color.HiRed("[!] invalid url %s: %s\n", target, err)
		return
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		color.HiRed("[!] failed to fetch %s: %s\n", target, err)
		return
	}
	defer resp.Body.Close()

	// Limit response body to 5MB to prevent memory issues
	limitedReader := io.LimitReader(resp.Body, 5*1024*1024)
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		color.HiRed("[!] failed to read response from %s: %s\n", target, err)
		return
	}

	responseString := cleanResponse(string(body))
	foundS3 := strings.Contains(responseString, "amazonaws.com")
	foundCF := strings.Contains(responseString, "cloudfront.net")

	if foundS3 {
		color.HiGreen("[*] S3 bucket found! %s\n", target)
		if writeMatches(responseString, target, awsRegexes, s3path, &s3pathMu) {
			found.Add(1)
		}
	}
	if foundCF {
		color.HiGreen("[*] Cloudfront found! %s\n", target)
		if writeMatches(responseString, target, cloudfrontRegexes, cfpath, &cfpathMu) {
			found.Add(1)
		}
	}
	if !foundS3 && !foundCF {
		color.HiRed("[*] Moving to Next Site\n")
	}
}

// writeMatches normalises and de-duplicates matches, then appends them to the
// CSV output. Returns true when at least one match was written.
func writeMatches(responseString, sourceURL string, regexes []*regexp.Regexp, outputPath string, mu *sync.Mutex) bool {
	seen := make(map[string]bool)
	wrote := false

	for _, re := range regexes {
		match := re.FindString(responseString)
		if match == "" {
			continue
		}

		// Normalise URL to always use https://
		normalized := match
		if strings.HasPrefix(match, "//") {
			normalized = "https:" + match
		} else if strings.HasPrefix(match, "http://") {
			normalized = strings.Replace(match, "http://", "https://", 1)
		}

		// Deduplicate - only write unique normalized URLs
		if !seen[normalized] {
			seen[normalized] = true
			color.HiGreen("[*] Logging... %s\n", normalized)
			appendCSV(outputPath, sourceURL, normalized, mu)
			wrote = true
		}
	}
	return wrote
}

func newClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: *insecure}, // #nosec G402 - intentional, controlled by -insecure
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       30 * time.Second,
	}
	return &http.Client{Transport: tr, Timeout: *timeout}
}

func main() {
	// Parse the flags provided
	flag.Parse()
	color.HiGreen("[*]                                 [*]\n")
	color.HiGreen("[*] AWS Scanner - By @random_robbie [*]\n")
	color.HiGreen("[*]                                 [*]\n")
	color.HiGreen("[*]    S3 Buckets for the Win       [*]\n")
	color.HiGreen("[*]                                 [*]\n\n")

	if *list == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *concurr < 1 {
		*concurr = 1
	}

	// Read list line by line
	lines, err := readLines(*list)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client := newClient()
	urls := make(chan string)
	var wg sync.WaitGroup
	var found atomic.Int64

	wg.Add(*concurr)
	for range *concurr {
		go func() {
			defer wg.Done()
			for raw := range urls {
				select {
				case <-ctx.Done():
					return
				default:
				}
				scanURL(client, raw, &found)
			}
		}()
	}

	for _, line := range lines {
		select {
		case urls <- line:
		case <-ctx.Done():
		}
		if ctx.Err() != nil {
			break
		}
	}
	close(urls)
	wg.Wait()

	color.HiGreen("[*] done: %d findings — see %s and %s\n", found.Load(), s3path, cfpath)
	if found.Load() > 0 {
		// Non-zero exit when findings exist, so pipelines can react.
		os.Exit(1)
	}
}
