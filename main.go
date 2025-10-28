package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	tr      = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	list    = flag.String("list", "list.txt", "Default: list.txt")
	s3path  = "s3-bucket.csv"
	cfpath  = "cloudfront.csv"
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
)

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return err != nil
}

// Read a whole file into the memory and store it as array of lines
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

// write file of s3bucket
func writes3bucket(s3bucket, url, s3path string) {
	fileHandle, err := os.OpenFile(s3path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if isError(err) {
		return
	}
	defer fileHandle.Close()

	writer := bufio.NewWriter(fileHandle)
	stringtolog := fmt.Sprintf("%s,%s", url, s3bucket)
	fmt.Fprintln(writer, stringtolog)
	writer.Flush()
}

// write file of cloudfront
func writecf(cf, url, cfpath string) {
	fileHandle, err := os.OpenFile(cfpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if isError(err) {
		return
	}
	defer fileHandle.Close()

	writer := bufio.NewWriter(fileHandle)
	stringtolog := fmt.Sprintf("%s,%s", url, cf)
	fmt.Fprintln(writer, stringtolog)
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

func runRegexCheck(responseString, url string, patterns []string, outputPath string, writeFunc func(string, string, string)) {
	seen := make(map[string]bool)

	for _, pattern := range patterns {
		re, _ := regexp.Compile(pattern)
		match := re.FindString(responseString)

		if match != "" {
			// Normalize URL to always use https://
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
				writeFunc(normalized, url, outputPath)
			}
		}
	}
}

func main() {
	// Parse the flags provided
	flag.Parse()
	color.HiGreen("[*]                                 [*]\n")
	color.HiGreen("[*] AWS Scanner - By @random_robbie [*]\n")
	color.HiGreen("[*]                                 [*]\n")
	color.HiGreen("[*]                                 [*]\n")
	color.HiGreen("[*]    S3 Buckets for the Win       [*]\n")
	color.HiGreen("[*]                                 [*]\n")
	color.HiGreen("[*]       Version 1.0               [*]\n\n\n\n")

	//Check the server is not blank or empty
	if *list == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	//read list line by line
	lines, err := readLines(*list)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	for _, line := range lines {

		//Print Scanning Website
		url := normalizeURL(line)
		attempt := fmt.Sprintf("[*] Scanning Website %s\n", url)
		color.HiWhite(attempt)
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
		// Do not verify certificates
		client := &http.Client{
			Transport: tr,
		}

		// Create HEAD request with random user agent.
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0")

		if resp, err := client.Do(req); err == nil {

			defer resp.Body.Close()
			// Limit response body to 5MB to prevent memory issues
			limitedReader := io.LimitReader(resp.Body, 5*1024*1024)
			body, err := io.ReadAll(limitedReader)
			if err != nil {
				panic(err)
			}
			// Clean response to remove messy JS
			responseString := cleanResponse(string(body))
			//If amazonaws.com is found signal bucket found
			if strings.Contains(responseString, "amazonaws.com") == true {
				color.HiGreen("[*] S3 bucket found!\n")
				runRegexCheck(responseString, url, awsurls, s3path, writes3bucket)
			}
			//find Cloudfront distribution end point.
			if strings.Contains(responseString, "cloudfront.net") == true {
				color.HiGreen("[*] Cloudfront found!\n")
				runRegexCheck(responseString, url, cloudfronturls, cfpath, writecf)

			} else {
				color.HiRed("[*] Moving to Next Site\n")
			}
		}
	}

}
