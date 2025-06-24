package crawler

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ameer005/meowl/internals/repository"
	"github.com/ameer005/meowl/pkg/queue"
	"go.mongodb.org/mongo-driver/mongo"
)

type Crawler struct {
	queue         *queue.Queue[string]
	seen          map[string]struct{}
	logger        *slog.Logger
	websiteRepo   *repository.WebsiteRepo
	Wg            sync.WaitGroup
	mu            sync.Mutex
	activeWorkers int
}

func New(urls []string, logger *slog.Logger, db *mongo.Client) *Crawler {

	q := queue.NewQueue[string]()
	for _, url := range urls {
		q.Enqueue(url)
	}
	return &Crawler{
		queue:         q,
		seen:          make(map[string]struct{}),
		logger:        logger,
		websiteRepo:   repository.NewWebsiteRepo(db.Database("spider")),
		Wg:            sync.WaitGroup{},
		activeWorkers: 0,
	}
}

func (t *Crawler) Start(ctx context.Context) {
	defer func() {
		t.Wg.Done()
	}()

	// bfs loop
	for {

		t.mu.Lock()

		if t.queue.IsEmpty() {
			if t.activeWorkers == 0 {
				t.mu.Unlock()
				t.logger.Info("All links have been scraped")
				break
			}

			// Improve this sleeping time
			time.Sleep(100 * time.Millisecond)
			t.mu.Unlock()
			continue
		}

		// extracting URL
		url, ok := t.queue.Dequeue()
		if !ok {
			t.mu.Unlock()
			t.logger.Warn("Queue is empty")
			break
		}

		if _, isSeen := t.seen[url]; isSeen {
			t.mu.Unlock()
			continue
		}

		t.seen[url] = struct{}{}
		t.activeWorkers++
		t.mu.Unlock()

		// processing url
		htmlStr, statusCode, contentType, err := t.fetchHTML(url)
		if err != nil {
			t.logger.Error("Fetch HTML error", slog.String("url", url), slog.String("error", err.Error()))

			t.releaseWorker()
			continue
		}

		// Handling http errors
		if statusCode == 400 {
			t.logger.Warn("400 error")
			continue
		}

		if statusCode == 401 {
			t.logger.Warn("Authentication Error! skipping...")
			continue

		}

		if statusCode == 403 {
			t.logger.Warn("Authorization error! skipping...")
			continue
		}

		if statusCode != 200 {
			continue
		}

		if contentType == "" {
		}

		if statusCode != 200 {
			t.releaseWorker()
			continue
		}

		// Parsing HTML
		t.logger.Info("Fetched",
			slog.String("url", url),
			slog.Int("status_code", statusCode),
		)

		reader := strings.NewReader(htmlStr)
		websiteData, err := extractContent(reader, url)
		if err != nil {
			t.logger.Error(err.Error())
		}

		// updating queue
		t.mu.Lock()
		for _, newURL := range websiteData.Outlinks {
			if _, isSeen := t.seen[newURL]; isSeen {
				continue
			}

			t.queue.Enqueue(newURL)
		}
		t.mu.Unlock()

		t.logger.Info("Fetched", "url", url)

		fetchedWebsite, err := t.websiteRepo.GetByURL(ctx, url)

		if err != nil {
			t.logger.Warn("fetch website error", "error", err)
		}

		// TODO: maybe move this logic to upwards
		// data exist
		if fetchedWebsite != nil {
			t.releaseWorker()
			continue
		}

		err = t.websiteRepo.AddWebsite(ctx, websiteData)
		if err != nil {
			t.logger.Error("Inser website to db error", "error", err)
		}

		t.releaseWorker()

	}

	t.logger.Info("Crawler existed", "active workers", t.activeWorkers)

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

		// bodyBytes = append(bodyBytes,  zzz)
	}

	return string(bodyBytes), r.StatusCode, contentType, nil
}

func (c *Crawler) releaseWorker() {
	c.mu.Lock()
	c.activeWorkers--
	c.mu.Unlock()
}
