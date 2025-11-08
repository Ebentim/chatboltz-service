package scraper

// Data shapes for the scraper service
type ScrapeOptions struct {
	Trace    bool     `json:"trace"`     // whether to follow links within same domain
	Exclude  []string `json:"exclude"`   // list of URL patterns to exclude (supports simple wildcards, e.g. "https://site.com/login*")
	MaxPages int      `json:"max_pages"` // maximum pages to fetch when tracing
}

type Section struct {
	Tag  string `json:"tag"`
	Text string `json:"text"`
}

type PageData struct {
	URL      string    `json:"url"`
	Title    string    `json:"title,omitempty"`
	Sections []Section `json:"sections,omitempty"`
	Links    []string  `json:"links,omitempty"`
	Images   []string  `json:"images,omitempty"`
}

type ScrapeResult struct {
	Pages []PageData `json:"pages"`
}
