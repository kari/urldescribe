package urldescribe

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	html "golang.org/x/net/html"
)

// DescribeURL takes an URL and returns its description (ie. its <title> tag)
func DescribeURL(rawurl string) string {

	url, err := parseURL(rawurl)
	if err != nil {
		return "" // Can't parse input
	}

	resp, err := getPage(url.String())
	if err != nil {
		return "" // error fetching page
	}
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err.Error() // Error parsing HTML
	}

	title := strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(extractTitle(doc)), "\n", " "), "  ", " ") // remove newlines and convert ensuing double spaces to single

	// IRC has limit of about 510 bytes for message length
	if len([]rune(title)) > 140 {
		title = string([]rune(title)[0:139]) + "…"
	}

	return title
}

// getPage tries to GET a url and returns a http.Response if the
// response has status 200 OK and it is a text/html document
func getPage(url string) (*http.Response, error) {

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err // Creating request failed
	}
	req.Header.Set("Accept", "text/html")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err // Request failed
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Request responded with status %d", resp.StatusCode) // Request didn't return 200 OK
	}
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html") {
		// fmt.Println(resp.Header.Get("Content-Type"))
		return nil, fmt.Errorf("Response is %s, not text/html", resp.Header.Get("Content-Type")) // Not an HTML document
	}
	if resp.ContentLength > 2*1024*1024 {
		// Response body is larger than 2 MB
		return nil, fmt.Errorf("Response body is too large (%d)", resp.ContentLength) // see: https://github.com/cloudfoundry/bytefmt
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

	var f func(*html.Node) string

	f = func(n *html.Node) string {
		if n.Type == html.ElementNode && n.Data == "title" {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					return c.Data
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if f(c) != "" {
				return f(c)
			}
		}

		return ""
	}

	return f(tree)
}

// parseURL adds limitations to url.ParseRequestURI that
// the url needs to have a host and it can't be localhost
func parseURL(rawurl string) (*url.URL, error) {

	url, err := url.ParseRequestURI(rawurl)

	if err != nil {
		return nil, err
	}
	if strings.ToLower(url.Host) == "localhost" {
		return nil, errors.New("URL cannot point to localhost")
	}
	if url.Host == "" {
		return nil, errors.New("URL needs to have a host")
	}

	return url, nil
}
