package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// These are live integration-style tests that hit the real host
// They run only when the environment variable RUN_LIVE_SCRAPER_TESTS=1 is set.
// This prevents accidental network calls during CI or local quick tests.

func shouldRunLiveTests() bool {
	return os.Getenv("RUN_LIVE_SCRAPER_TESTS") == "1"
}

// helper that writes JSON test output to tmp/scraper-tests/<name>.json
func writeTestResult(filename string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		return err
	}
	return os.WriteFile(filename, b, 0o644)
}

// writeTestResultText writes a human-readable text summary of the scrape result
func writeTestResultText(filename string, v any) error {
	var sb strings.Builder

	switch t := v.(type) {
	case ScrapeResult:
		for i, p := range t.Pages {
			sb.WriteString(fmt.Sprintf("Page %d:\n", i+1))
			sb.WriteString(fmt.Sprintf("URL: %s\n", p.URL))
			if p.Title != "" {
				sb.WriteString(fmt.Sprintf("Title: %s\n", p.Title))
			}
			if len(p.Sections) > 0 {
				sb.WriteString("Sections:\n")
				for _, s := range p.Sections {
					text := s.Text
					if len(text) > 300 {
						text = text[:300] + "..."
					}
					sb.WriteString(fmt.Sprintf("  - [%s] %s\n", s.Tag, text))
				}
			}
			if len(p.Links) > 0 {
				sb.WriteString("Links:\n")
				for _, l := range p.Links {
					sb.WriteString(fmt.Sprintf("  - %s\n", l))
				}
			}
			if len(p.Images) > 0 {
				sb.WriteString("Images:\n")
				for _, im := range p.Images {
					sb.WriteString(fmt.Sprintf("  - %s\n", im))
				}
			}
			sb.WriteString("\n---\n\n")
		}
	case []PageData:
		for i, p := range t {
			sb.WriteString(fmt.Sprintf("Page %d:\n", i+1))
			sb.WriteString(fmt.Sprintf("URL: %s\n", p.URL))
			if p.Title != "" {
				sb.WriteString(fmt.Sprintf("Title: %s\n", p.Title))
			}
			if len(p.Sections) > 0 {
				sb.WriteString("Sections:\n")
				for _, s := range p.Sections {
					text := s.Text
					if len(text) > 300 {
						text = text[:300] + "..."
					}
					sb.WriteString(fmt.Sprintf("  - [%s] %s\n", s.Tag, text))
				}
			}
			if len(p.Links) > 0 {
				sb.WriteString("Links:\n")
				for _, l := range p.Links {
					sb.WriteString(fmt.Sprintf("  - %s\n", l))
				}
			}
			if len(p.Images) > 0 {
				sb.WriteString("Images:\n")
				for _, im := range p.Images {
					sb.WriteString(fmt.Sprintf("  - %s\n", im))
				}
			}
			sb.WriteString("\n---\n\n")
		}
	case map[string]string:
		for k, vv := range t {
			sb.WriteString(fmt.Sprintf("%s: %s\n", k, vv))
		}
	default:
		// fallback: try JSON pretty print
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			sb.WriteString(fmt.Sprintf("(unable to format result): %v\n", err))
		} else {
			sb.WriteString(string(b))
		}
	}

	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		return err
	}
	return os.WriteFile(filename, []byte(sb.String()), 0o644)
}

func TestScrapeDomain_AvoidAuthPages(t *testing.T) {
	outDir := filepath.Join("tmp", "scraper-tests")
	_ = os.MkdirAll(outDir, 0o755)
	if !shouldRunLiveTests() {
		// write skip file
		_ = writeTestResult(filepath.Join(outDir, "TestScrapeDomain_AvoidAuthPages.json"), map[string]string{"status": "skipped", "reason": "RUN_LIVE_SCRAPER_TESTS!=1"})
		t.Skip("Skipping live scraper tests. Set RUN_LIVE_SCRAPER_TESTS=1 to enable.")
	}

	root := "https://www.alpinesbolt.com"
	svc := NewService(&http.Client{Timeout: 20 * time.Second})

	// Exclude common auth paths to avoid login/register pages
	excludes := []string{
		root + "/login",
		root + "/signin",
		root + "/register",
		root + "/signup",
		root + "/account",
	}

	res, err := svc.Scrape(context.Background(), root, ScrapeOptions{Trace: true, Exclude: excludes, MaxPages: 30})
	if err != nil {
		_ = writeTestResult(filepath.Join(outDir, "TestScrapeDomain_AvoidAuthPages.json"), map[string]string{"status": "error", "error": err.Error()})
		t.Fatalf("live scrape failed: %v", err)
	}

	if len(res.Pages) == 0 {
		t.Fatalf("expected at least one page from %s", root)
	}

	for _, p := range res.Pages {
		if !strings.Contains(p.URL, "alpinesbolt.com") {
			_ = writeTestResult(filepath.Join(outDir, "TestScrapeDomain_AvoidAuthPages.json"), map[string]string{"status": "error", "error": "visited outside domain: " + p.URL})
			t.Fatalf("scraper visited outside domain: %s", p.URL)
		}
		for _, ex := range excludes {
			if p.URL == ex {
				_ = writeTestResult(filepath.Join(outDir, "TestScrapeDomain_AvoidAuthPages.json"), map[string]string{"status": "error", "error": "visited excluded URL: " + ex})
				t.Fatalf("scraper visited excluded URL: %s", ex)
			}
		}
	}

	// write successful result (JSON + TXT)
	jsonFile := filepath.Join(outDir, "TestScrapeDomain_AvoidAuthPages.json")
	txtFile := filepath.Join(outDir, "TestScrapeDomain_AvoidAuthPages.txt")
	_ = writeTestResult(jsonFile, res)
	_ = writeTestResultText(txtFile, res)
}

func TestScrapeProducts_Only(t *testing.T) {
	outDir := filepath.Join("tmp", "scraper-tests")
	_ = os.MkdirAll(outDir, 0o755)
	if !shouldRunLiveTests() {
		_ = writeTestResult(filepath.Join(outDir, "TestScrapeProducts_Only.json"), map[string]string{"status": "skipped", "reason": "RUN_LIVE_SCRAPER_TESTS!=1"})
		t.Skip("Skipping live scraper tests. Set RUN_LIVE_SCRAPER_TESTS=1 to enable.")
	}

	root := "https://www.alpinesbolt.com"
	svc := NewService(&http.Client{Timeout: 20 * time.Second})

	// Candidate product-listing paths to probe; we'll pick the first that exists
	candidates := []string{"/products", "/shop", "/collections", "/catalog"}
	client := &http.Client{Timeout: 8 * time.Second}
	var start string
	for _, p := range candidates {
		u := root + p
		req, _ := http.NewRequest(http.MethodHead, u, nil)
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 400 {
			start = u
			break
		}
	}
	if start == "" {
		_ = writeTestResult(filepath.Join(outDir, "TestScrapeProducts_Only.json"), map[string]string{"status": "skipped", "reason": "no product listing path found"})
		t.Skip("No product listing path found on host; skipping product-only test")
	}

	res, err := svc.Scrape(context.Background(), start, ScrapeOptions{Trace: true, MaxPages: 50})
	if err != nil {
		_ = writeTestResult(filepath.Join(outDir, "TestScrapeProducts_Only.json"), map[string]string{"status": "error", "error": err.Error()})
		t.Fatalf("scrape products failed: %v", err)
	}

	// Filter pages that look like product pages by URL or by content hints
	products := []PageData{}
	for _, p := range res.Pages {
		if strings.Contains(strings.ToLower(p.URL), "product") || strings.Contains(strings.ToLower(p.URL), "/p/") {
			products = append(products, p)
			continue
		}
		// also check page text for "product" keyword
		for _, s := range p.Sections {
			if strings.Contains(strings.ToLower(s.Text), "product") {
				products = append(products, p)
				break
			}
		}
	}

	if len(products) == 0 {
		_ = writeTestResult(filepath.Join(outDir, "TestScrapeProducts_Only.json"), res)
		_ = writeTestResultText(filepath.Join(outDir, "TestScrapeProducts_Only.txt"), res)
		t.Fatalf("expected to find at least one product page starting from %s; found none", start)
	}

	// write successful product subset (JSON + TXT)
	_ = writeTestResult(filepath.Join(outDir, "TestScrapeProducts_Only.json"), products)
	_ = writeTestResultText(filepath.Join(outDir, "TestScrapeProducts_Only.txt"), products)
}

func TestScrapeDomain_WholeSite(t *testing.T) {
	outDir := filepath.Join("tmp", "scraper-tests")
	_ = os.MkdirAll(outDir, 0o755)
	if !shouldRunLiveTests() {
		_ = writeTestResult(filepath.Join(outDir, "TestScrapeDomain_WholeSite.json"), map[string]string{"status": "skipped", "reason": "RUN_LIVE_SCRAPER_TESTS!=1"})
		t.Skip("Skipping live scraper tests. Set RUN_LIVE_SCRAPER_TESTS=1 to enable.")
	}

	root := "https://www.alpinesbolt.com"
	svc := NewService(&http.Client{Timeout: 30 * time.Second})

	res, err := svc.Scrape(context.Background(), root, ScrapeOptions{Trace: true, MaxPages: 100})
	if err != nil {
		_ = writeTestResult(filepath.Join(outDir, "TestScrapeDomain_WholeSite.json"), map[string]string{"status": "error", "error": err.Error()})
		t.Fatalf("full-site scrape failed: %v", err)
	}

	if len(res.Pages) == 0 {
		_ = writeTestResult(filepath.Join(outDir, "TestScrapeDomain_WholeSite.json"), map[string]string{"status": "error", "error": "no pages scraped"})
		_ = writeTestResultText(filepath.Join(outDir, "TestScrapeDomain_WholeSite.txt"), map[string]string{"status": "error", "error": "no pages scraped"})
		t.Fatalf("expected to scrape pages from %s", root)
	}
	for _, p := range res.Pages {
		if !strings.Contains(p.URL, "alpinesbolt.com") {
			_ = writeTestResult(filepath.Join(outDir, "TestScrapeDomain_WholeSite.json"), map[string]string{"status": "error", "error": "visited outside domain: " + p.URL})
			_ = writeTestResultText(filepath.Join(outDir, "TestScrapeDomain_WholeSite.txt"), map[string]string{"status": "error", "error": "visited outside domain: " + p.URL})
			t.Fatalf("scraper visited outside domain during full-site run: %s", p.URL)
		}
	}

	_ = writeTestResult(filepath.Join(outDir, "TestScrapeDomain_WholeSite.json"), res)
	_ = writeTestResultText(filepath.Join(outDir, "TestScrapeDomain_WholeSite.txt"), res)
}
