package crawler

import (
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/ameer005/meowl/internals/models"
	"github.com/lib/pq"
	"golang.org/x/net/html"
)

func extractContent(reader io.Reader, domain string) (*models.Website, error) {
	doc, err := html.Parse(reader)
	website := models.Website{
		Url:           domain,
		Title:         "",
		Headings:      pq.StringArray{},
		Content:       "",
		InternalLinks: pq.StringArray{},
		ExternalLinks: pq.StringArray{},
		Images:        pq.StringArray{},
		Description:   "",
	}

	if err != nil {
		return &website, fmt.Errorf("Parser:Parsing html error: %v", err)
	}

	parsedDomain, err := url.ParseRequestURI(domain)
	if err != nil {
		return &website, err
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

				processedURL := extractURL(attr.Val, domain)
				if processedURL != "" {
					// separating internal and external Links
					parsedURL, err := url.ParseRequestURI(processedURL)
					if err != nil {
						continue
					}

					domainHost := parsedDomain.Hostname()
					linkHost := parsedURL.Hostname()

					if linkHost == domainHost {
						website.InternalLinks = append(website.InternalLinks, processedURL)
					} else {
						website.ExternalLinks = append(website.ExternalLinks, processedURL)
					}
				}
			}
		}

		// parsing headings
		if n.Type == html.ElementNode && n.Data == "title" {
			if n.FirstChild != nil {
				title := strings.TrimSpace(n.FirstChild.Data)
				if title != "" {
					website.Title = title
				}

			}
			website.Title = n.FirstChild.Data
		}

		// parsing meta description
		if n.Type == html.ElementNode && n.Data == "meta" {
			isDescription := false
			for _, attr := range n.Attr {
				if attr.Val == "description" {
					isDescription = true
				}
			}

			if isDescription {
				for _, attr := range n.Attr {
					if attr.Key == "content" {
						if attr.Val != "" {
							fmt.Println(attr.Val)
						}
					}
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
			website.Headings = append(website.Headings, headingStr)

		}

		//parsing text content
		if n.Type == html.TextNode && n.Parent != nil {
			if n.Parent.Data != "script" && n.Parent.Data != "style" &&
				n.Parent.Data != "h1" && n.Parent.Data != "h2" && n.Parent.Data != "h3" &&
				n.Parent.Data != "h4" && n.Parent.Data != "h5" && n.Parent.Data != "h6" {

				text := strings.TrimSpace(n.Data)
				if text != "" {
					website.Content += text
				}
			}
		}

		// images
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, img := range n.Attr {
				if img.Key == "src" && !strings.HasPrefix(img.Val, "data:") {
					imageSrc := extractURL(img.Val, domain)
					website.Images = append(website.Images, imageSrc)
				}
			}

		}

		f(n.FirstChild)
		f(n.NextSibling)
	}

	f(doc)

	website.CrawledAt = time.Now()
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
