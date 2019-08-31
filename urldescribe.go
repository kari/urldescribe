package urldescribe

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	html "golang.org/x/net/html"
)

// DescribeURL takes an URL and returns its description
func DescribeURL(rawurl string) string {

	url, err := parseURL(rawurl)
	if err != nil {
		return "" // Can't parse input
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return "" // Creating request failed
	}
	req.Header.Set("Accept", "text/html")

	resp, err := client.Do(req)
	if err != nil {
		return "" // Request failed
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "" // Request didn't return 200 OK
	}
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html") {
		// fmt.Println(resp.Header.Get("Content-Type"))
		return "" // Not an HTML document
	}
	if resp.ContentLength > 2*1024*1024 {
		return "" // Response body is larger than 2 MB
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err.Error() // Error parsing HTML
	}

	title := strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(extractTitle(doc)), "\n", " "), "  ", " ") // remove newlines and convert ensuing double spaces to single

	// IRC has limit of about 510 bytes for message length
	if len([]rune(title)) > 140 {
		title = string([]rune(title)[0:139]) + "..."
	}

	return title
}

func extractTitle(tree *html.Node) string {
	var f func(*html.Node) string
	f = func(n *html.Node) string {
		if n.Type == html.ElementNode && n.Data == "title" {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					// fmt.Println(c.Data)
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
