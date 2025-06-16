package crawler

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Crawler struct {
	visited map[string]struct{}
	queue   []string
	logger  *slog.Logger
}

type Website struct {
	url       string
	content   string
	outlinks  []string
	backlinks []string
}

func New(urls []string, logger *slog.Logger) *Crawler {
	return &Crawler{
		queue:   urls,
		visited: make(map[string]struct{}),
		logger:  logger,
	}
}

func (t *Crawler) Start() {
	for _, url := range t.queue {
		t.parse(url)
	}

}

func (t *Crawler) parse(url string) {
	fmt.Printf("Fetching %s", url)

	r, err := http.Get(url)

	if err != nil {
		t.logger.Error("Fetching url error")
		return
	}

	contentType := r.Header.Get("Content-Type")

	fmt.Println("content type is: ", contentType)

	bodyBytes, err := io.ReadAll(r.Body)
	bodyString := string(bodyBytes)

	fmt.Println(bodyString)

	fmt.Printf("\n%+v", r)

}
