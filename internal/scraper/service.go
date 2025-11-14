package scraper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Service represents the scraper service
type Service struct {
	client *http.Client
}

// NewService creates a scraper service. If client is nil a default client with timeout is used.
func NewService(client *http.Client) *Service {
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	return &Service{client: client}
}

// Scrape fetches the provided URL and (optionally) traces internal links up to MaxPages.
func (s *Service) Scrape(ctx context.Context, rawurl string, opts ScrapeOptions) (ScrapeResult, error) {
	if rawurl == "" {
		return ScrapeResult{}, errors.New("url is required")
	}

	parsedRoot, err := url.Parse(rawurl)
	if err != nil {
		return ScrapeResult{}, err
	}

	// exclude patterns (supports wildcards). Keep as slice and match per-URL.
	excludePatterns := opts.Exclude

	maxPages := opts.MaxPages
	if maxPages <= 0 {
		maxPages = 20
	}

	pages := []PageData{}
	visited := map[string]struct{}{}

	type queueItem struct {
		u     *url.URL
		depth int
	}

	q := []queueItem{{u: parsedRoot, depth: 0}}

	for len(q) > 0 && len(pages) < maxPages {
		item := q[0]
		q = q[1:]
		u := item.u
		norm := u.String()

		if _, ok := visited[norm]; ok {
			continue
		}
		if isExcluded(norm, excludePatterns) {
			visited[norm] = struct{}{}
			continue
		}

		// fetch page
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, norm, nil)
		if err != nil {
			// skip this URL
			visited[norm] = struct{}{}
			continue
		}
		resp, err := s.client.Do(req)
		if err != nil {
			visited[norm] = struct{}{}
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			visited[norm] = struct{}{}
			continue
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
		if err != nil {
			visited[norm] = struct{}{}
			continue
		}

		pd := PageData{URL: norm}
		pd.Title = strings.TrimSpace(doc.Find("title").Text())

		// extract sections (headings and paragraphs)
		sections := []Section{}
		for i := 1; i <= 6; i++ {
			tag := fmt.Sprintf("h%d", i)
			doc.Find(tag).Each(func(_ int, s *goquery.Selection) {
				text := strings.TrimSpace(s.Text())
				if text != "" {
					sections = append(sections, Section{Tag: tag, Text: text})
				}
			})
		}
		// paragraphs
		doc.Find("p").Each(func(_ int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if text != "" {
				sections = append(sections, Section{Tag: "p", Text: text})
			}
		})
		pd.Sections = sections

		// links
		links := []string{}
		doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			if href == "" {
				return
			}
			resolved := resolveURL(u, href)
			if resolved != "" {
				links = append(links, resolved)
			}
		})
		pd.Links = uniqueStrings(links)

		// images
		images := []string{}
		doc.Find("img[src]").Each(func(_ int, s *goquery.Selection) {
			src, _ := s.Attr("src")
			if src == "" {
				return
			}
			resolved := resolveURL(u, src)
			if resolved != "" {
				images = append(images, resolved)
			}
		})
		pd.Images = uniqueStrings(images)

		pages = append(pages, pd)
		visited[norm] = struct{}{}

		// tracing: enqueue links within same host
		if opts.Trace && len(pages) < maxPages {
			for _, l := range pd.Links {
				if _, ok := visited[l]; ok {
					continue
				}
				if isExcluded(l, excludePatterns) {
					continue
				}
				parsed, err := url.Parse(l)
				if err != nil {
					continue
				}
				if sameHost(parsedRoot, parsed) {
					q = append(q, queueItem{u: parsed, depth: item.depth + 1})
				}
			}
		}
	}

	return ScrapeResult{Pages: pages}, nil
}

func resolveURL(base *url.URL, href string) string {
	href = strings.TrimSpace(href)
	if href == "" {
		return ""
	}
	parsed, err := url.Parse(href)
	if err != nil {
		return ""
	}
	resolved := base.ResolveReference(parsed)
	return resolved.String()
}

func sameHost(a, b *url.URL) bool {
	return a.Hostname() == b.Hostname()
}

func uniqueStrings(in []string) []string {
	m := map[string]struct{}{}
	out := []string{}
	for _, s := range in {
		if _, ok := m[s]; !ok {
			m[s] = struct{}{}
			out = append(out, s)
		}
	}
	return out
}

// isExcluded returns true if the given URL matches any of the exclude patterns.
// Patterns support simple shell-style wildcards (*, ?). Matching is done against the full URL string.
func isExcluded(u string, patterns []string) bool {
	if len(patterns) == 0 {
		return false
	}
	for _, p := range patterns {
		if p == "" {
			continue
		}
		// If pattern contains wildcard characters, use path.Match
		if strings.ContainsAny(p, "*?[]") {
			ok, err := path.Match(p, u)
			if err == nil && ok {
				return true
			}
			// also try matching against the URL's path portion
			if parsed, err := url.Parse(u); err == nil {
				if ok, _ := path.Match(p, parsed.Path); ok {
					return true
				}
			}
		} else {
			// exact match
			if p == u {
				return true
			}
		}
	}
	return false
}
