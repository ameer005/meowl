package crawler

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func extractContent(reader io.Reader, domain string) ([]string, error) {
	fmt.Println("running extract content")
	urls := []string{}
	doc, err := html.Parse(reader)

	if err != nil {
		return urls, fmt.Errorf("Parsing html error: %v", err)
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n == nil {
			return
		}

		// links
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key != "href" {
					continue
				}

				url := extractURL(attr.Val, domain)
				if url != "" {
					urls = append(urls, url)
				}
			}
		}

		// text content
		if n.Type == html.TextNode {
			if n.Parent != nil && (n.Parent.Data == "script" || n.Parent.Data == "style") {
				return
			}

		}

		f(n.FirstChild)
		f(n.NextSibling)

	}

	f(doc)

	return urls, nil

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
