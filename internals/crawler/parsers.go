package crawler

import (
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Website struct {
	url string
	// TODO use string builder for perforance
	content   string
	title     []string
	headings  string
	outlinks  []string
	images    []string
	crawledAt time.Time
}

func extractContent(reader io.Reader, domain string) (*Website, error) {
	fmt.Println("running extract content")
	doc, err := html.Parse(reader)
	website := Website{}

	if err != nil {
		return &website, fmt.Errorf("Parsing html error: %v", err)
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n == nil {
			return
		}

		//parsing links
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key != "href" {
					continue
				}

				url := extractURL(attr.Val, domain)
				if url != "" {
					website.outlinks = append(website.outlinks, url)
				}
			}
		}

		//parsing headings
		isHeading := (n.Data == "h1" || n.Data == "h2" || n.Data == "h3" || n.Data == "h4" || n.Data == "h5" || n.Data == "h6")

		if n.Type == html.ElementNode && isHeading {
			var headingStr string

			var extractText func(*html.Node)
			extractText = func(c *html.Node) {
				if c.Type == html.TextNode {
					headingStr += c.Data
				}

				for child := c.FirstChild; child != nil; child = child.NextSibling {
					extractText(child)
				}
			}

			extractText(n)
			website.headings = headingStr

		}

		//parsing text content
		if n.Type == html.TextNode && n.Parent != nil {
			if n.Parent.Data != "script" && n.Parent.Data != "style" &&
				n.Parent.Data != "h1" && n.Parent.Data != "h2" && n.Parent.Data != "h3" &&
				n.Parent.Data != "h4" && n.Parent.Data != "h5" && n.Parent.Data != "h6" {

				text := strings.TrimSpace(n.Data)
				if text != "" {
					website.content += text
				}
			}
		}

		f(n.FirstChild)
		f(n.NextSibling)

	}

	f(doc)

	return &website, nil

}

func extractURL(val string, domain string) string {
	val = strings.TrimSpace(val)

	if val == "" || val == "#" {
		return ""
	}

	if strings.HasPrefix(val, "//") {
		// Protocol-relative URL (e.g. //example.com/path)
		val = "https:" + val
	} else if strings.HasPrefix(val, "/") {
		// Relative path to current domain
		if len(val) == 1 {
			return ""
		}
		val = domain + val
	} else if !strings.HasPrefix(val, "http://") && !strings.HasPrefix(val, "https://") {
		return ""
	}

	parsed, err := url.ParseRequestURI(val)
	if err != nil {
		return ""
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ""
	}

	if parsed.Host == "" {
		return ""
	}

	return parsed.Scheme + "://" + parsed.Host + parsed.EscapedPath()
}
