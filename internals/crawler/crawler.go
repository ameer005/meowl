package crawler

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type Crawler struct {
	visited map[string]struct{}
	queue   []string
	logger  *slog.Logger
}

func New(urls []string, logger *slog.Logger) *Crawler {
	return &Crawler{
		queue:   urls,
		visited: make(map[string]struct{}),
		logger:  logger,
	}
}

func (t *Crawler) Start() {

	processCounter := 0

	for {
		if processCounter >= len(t.queue) {
			t.logger.Info("All links have been scraped")
			return
		}

		url := t.queue[processCounter]

		htmlStr, statusCode, contentType, err := t.fetchHTML(url)
		if err != nil {
			//TODO handle other errors properly
			t.logger.Error("Fetch HTML error", slog.String("error", err.Error()))
		}
		t.logger.Info("Fetched",
			slog.String("url", url),
			slog.Int("status_code", statusCode),
		)

		reader := strings.NewReader(htmlStr)
		websiteData, err := extractContent(reader, url)
		if err != nil {
			t.logger.Error(err.Error())
		}

		fmt.Printf("\n %v", websiteData.url)

		if statusCode == 403 {
		}

		if contentType == "" {
		}

		fmt.Println()
		processCounter += 1
	}

}

func (t *Crawler) fetchHTML(url string) (string, int, string, error) {
	fmt.Printf("Fetching %s \n", url)

	r, err := http.Get(url)

	if err != nil {
		return "", r.StatusCode, "", fmt.Errorf("fetching url  %v", err)
	}

	defer r.Body.Close()

	contentHeaders := strings.SplitN(r.Header.Get("Content-Type"), ";", 2)

	if len(contentHeaders) == 0 {
		return "", r.StatusCode, "", fmt.Errorf("Paring content : %v", strings.Join(contentHeaders, ";"))
	}

	contentType := contentHeaders[0]

	if contentType != "text/html" {
		return "", r.StatusCode, contentType, fmt.Errorf("Invalid content type: %v", contentType)
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return "", r.StatusCode, contentType, fmt.Errorf("Body reading error: %v", err)
	}

	return string(bodyBytes), r.StatusCode, contentType, nil
}
