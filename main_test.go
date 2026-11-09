package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	cases := map[string]string{
		"example.com":   "https://example.com",
		"http://a.b/c":  "http://a.b/c",
		"https://a.b/c": "https://a.b/c",
	}
	for in, want := range cases {
		if got := normalizeURL(in); got != want {
			t.Errorf("normalizeURL(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestCleanResponse(t *testing.T) {
	in := `<html><body onload="evil()">hi<script>steal()</script><p onclick="x()">ok</p></body></html>`
	got := cleanResponse(in)
	if strings.Contains(got, "steal()") || strings.Contains(got, "evil()") || strings.Contains(got, "onclick") {
		t.Errorf("scripts should be stripped, got %q", got)
	}
	if !strings.Contains(got, "ok") {
		t.Errorf("body content should survive, got %q", got)
	}
}

func TestWriteMatchesDeduplicatesAndWrites(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "s3-bucket.csv")

	body := `<img src="//s3.amazonaws.com/bucket-one"> <a href="https://s3.amazonaws.com/bucket-one">` +
		`<a href="https://bucket-two.s3.amazonaws.com/x">`
	wrote := writeMatches(body, "https://example.com", awsRegexes, out, new(sync.Mutex))
	if !wrote {
		t.Fatal("expected matches to be written")
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read output: %s", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 unique findings (deduplicated), got %d: %q", len(lines), lines)
	}
	// The path-style pattern logs the object URL; the host-style pattern logs
	// the leftmost https://*.amazonaws.com host form it finds.
	if !strings.Contains(string(data), "https://s3.amazonaws.com/bucket-one") ||
		!strings.Contains(string(data), ",https://s3.amazonaws.com\n") {
		t.Errorf("expected normalised https findings, got %q", string(data))
	}
}

func TestScanURLAgainstLiveServer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`<script>x()</script><img src="https://s3.amazonaws.com/leaky-bucket/img.jpg">`))
	}))
	defer srv.Close()

	dir := t.TempDir()
	s3path = filepath.Join(dir, "s3-bucket.csv")
	cfpath = filepath.Join(dir, "cloudfront.csv")

	var found atomic.Int64
	scanURL(newClient(), srv.URL, &found)
	if found.Load() != 1 {
		t.Errorf("expected 1 finding, got %d", found.Load())
	}
	data, _ := os.ReadFile(s3path)
	if !strings.Contains(string(data), "https://s3.amazonaws.com/leaky-bucket") ||
		!strings.Contains(string(data), ",https://s3.amazonaws.com\n") {
		t.Errorf("expected leaky bucket logged, got %q", string(data))
	}
}

func TestScanURLSafeServerNoFindings(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("no clouds here"))
	}))
	defer srv.Close()

	s3path = filepath.Join(t.TempDir(), "s3-bucket.csv")
	cfpath = filepath.Join(t.TempDir(), "cloudfront.csv")

	var found atomic.Int64
	scanURL(newClient(), srv.URL, &found)
	if found.Load() != 0 {
		t.Errorf("expected no findings, got %d", found.Load())
	}
}
