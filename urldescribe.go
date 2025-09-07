package urldescribe

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	html "golang.org/x/net/html"
)

var (
	ErrInvalidHost      = errors.New("URL needs to have a host")
	ErrLocalhost        = errors.New("URL cannot point to localhost")
	ErrResponseTooLarge = errors.New("response body is too large")
	ErrNotHTML          = errors.New("response is not text/html")
)

// Config holds the configuration for URL description
type Config struct {
	timeout     time.Duration
	maxLength   int
	maxBodySize int64
	httpClient  *http.Client
}

// Option is a functional option for configuring the URL describer
type Option func(*Config)

// WithTimeout sets the HTTP client timeout
func WithTimeout(d time.Duration) Option {
	return func(c *Config) {
		c.timeout = d
	}
}

// WithMaxLength sets the maximum length for returned descriptions
func WithMaxLength(n int) Option {
	return func(c *Config) {
		c.maxLength = n
	}
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		timeout:     10 * time.Second,
		maxLength:   140,             // IRC has limit of about 510 bytes for message length
		maxBodySize: 2 * 1024 * 1024, // 2MB
		httpClient:  &http.Client{},
	}
}

// DescribeURL takes an URL and returns its description (ie. its <title> tag)
func DescribeURL(ctx context.Context, rawurl string, opts ...Option) (string, error) {
	cfg := DefaultConfig()

	for _, opt := range opts {
		opt(cfg)
	}

	parsedURL, err := parseURL(rawurl)
	if err != nil {
		return "", err
	}

	resp, err := getPage(ctx, cfg, parsedURL.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	title := formatTitle(extractTitle(doc), cfg.maxLength)
	return title, nil
}

func formatTitle(title string, maxLength int) string {
	var sb strings.Builder
	title = strings.TrimSpace(title)
	title = strings.ReplaceAll(title, "\n", " ")
	title = strings.Join(strings.Fields(title), " ") // handles multiple spaces

	runes := []rune(title)
	if len(runes) > maxLength {
		sb.WriteString(string(runes[:maxLength-1]))
		sb.WriteRune('…')
		return sb.String()
	}
	return title
}

// getPage tries to GET a url and returns a http.Response if the
// response has status 200 OK and it is a text/html document
func getPage(ctx context.Context, cfg *Config, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	req.Header.Set("User-Agent", "URLDescribe/1.0")

	resp, err := cfg.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html") {
		resp.Body.Close()
		return nil, fmt.Errorf("%w: %s", ErrNotHTML, resp.Header.Get("Content-Type"))
	}

	if resp.ContentLength > cfg.maxBodySize {
		resp.Body.Close()
		return nil, fmt.Errorf("%w: %d bytes", ErrResponseTooLarge, resp.ContentLength)
	}

	return resp, nil
}

// extractTitle finds the first <title> tag from a html.Node tree
// and returns it
func extractTitle(tree *html.Node) string {

	// According to HTML5 spec, <title> tag is only allowed within
	// a <head> block, which in turn is only allowed to be the first
	// child of a <html> element.

	// So, a depth-first search of a HTML page will return a proper
	// <title> tag or an improper one.

	var title string
	var crawler func(*html.Node)

	crawler = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && title == "" {
			if c := n.FirstChild; c != nil && c.Type == html.TextNode {
				title = c.Data
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawler(c)
		}
	}

	crawler(tree)
	return title
}

// parseURL adds limitations to url.ParseRequestURI that
// the url needs to have a host and it can't be localhost
func parseURL(rawurl string) (*url.URL, error) {
	parsedURL, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return nil, err
	}

	if parsedURL.Host == "" {
		return nil, ErrInvalidHost
	}

	if strings.ToLower(parsedURL.Host) == "localhost" {
		return nil, ErrLocalhost
	}

	return parsedURL, nil
}
