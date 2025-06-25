package crawler

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ameer005/meowl/internals/repository"
	"github.com/ameer005/meowl/pkg/queue"
	"go.mongodb.org/mongo-driver/mongo"
)

type Crawler struct {
	queue       *queue.Queue[string]
	seen        map[string]struct{}
	logger      *slog.Logger
	websiteRepo *repository.WebsiteRepo
	Wg          sync.WaitGroup
	mu          sync.Mutex
}

func New(urls []string, logger *slog.Logger, db *mongo.Client) *Crawler {

	q := queue.NewQueue[string]()
	for _, url := range urls {
		q.Enqueue(url)
	}
	return &Crawler{
		queue:       q,
		seen:        make(map[string]struct{}),
		logger:      logger,
		websiteRepo: repository.NewWebsiteRepo(db.Database("spider")),
		Wg:          sync.WaitGroup{},
	}
}

func (t *Crawler) Start(ctx context.Context) {
	defer func() {
		t.Wg.Done()
	}()
	workerID := rand.Intn(10000)
	t.logger.Info("Worker started", slog.Int("worker_id", workerID))

	idleCount := 0
	maxCount := 20

	// bfs loop
	for {

		t.mu.Lock()

		if t.queue.IsEmpty() {
			t.mu.Unlock()
			idleCount++
			if idleCount >= maxCount {
				t.logger.Info("Worker exiting due to idle timeout",
					slog.Int("worker_id", workerID),
					slog.Int("idle_count", idleCount),
				)

				return
			}

			// Improve this sleeping time
			time.Sleep(100 * time.Millisecond)
			continue
		}

		idleCount = 0

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
		t.mu.Unlock()

		// processing url
		htmlStr, statusCode, contentType, err := t.fetchHTML(url)
		if err != nil {
			t.logger.Error("Fetch failed",
				slog.String("url", url),
				slog.Int("worker_id", workerID),
				slog.String("error", err.Error()),
				slog.Int("status_code", statusCode),
			)
			continue
		}

		t.logger.Info("Fetched",
			slog.String("url", url),
			slog.Int("worker_id", workerID),
			slog.Int("status_code", statusCode),
			slog.String("content_type", contentType),
		)

		if statusCode != 200 {
			t.logger.Warn("Non-200 status code",
				slog.String("url", url),
				slog.Int("status_code", statusCode),
			)
			continue
		}

		if !strings.HasPrefix(contentType, "text/html") {
			t.logger.Warn("Skipping non-HTML content",
				slog.String("url", url),
				slog.String("content_type", contentType),
			)
			continue
		}

		// Parsing HTML
		reader := strings.NewReader(htmlStr)
		websiteData, err := extractContent(reader, url)
		if err != nil {
			t.logger.Error(err.Error())
		}

		// updating queue
		t.mu.Lock()
		for _, newURL := range websiteData.ExternalLinks {
			if _, isSeen := t.seen[newURL]; isSeen {
				continue
			}

			t.queue.Enqueue(newURL)
		}

		for _, newURL := range websiteData.InternalLinks {
			if _, isSeen := t.seen[newURL]; isSeen {
				continue
			}

			t.queue.Enqueue(newURL)
		}
		t.mu.Unlock()

		existing, err := t.websiteRepo.GetByURL(ctx, url)
		if err != nil {
			t.logger.Warn("DB lookup failed", slog.String("url", url), slog.String("error", err.Error()))
		}

		if existing == nil {
			if err := t.websiteRepo.AddWebsite(ctx, websiteData); err != nil {
				t.logger.Error("Failed to insert website to DB", slog.String("url", url), slog.String("error", err.Error()))
			} else {
				t.logger.Info("Website inserted to DB", slog.String("url", url))
			}
		} else {
			t.logger.Debug("Website already exists in DB", slog.String("url", url))
		}

		delay := randomDelay(2, 5)
		t.logger.Debug("Sleeping after URL", slog.Int("worker_id", workerID), slog.Duration("sleep", delay))
	}

}

func (t *Crawler) fetchHTML(url string) (string, int, string, error) {

	r, err := http.Get(url)

	if err != nil {

		return "", 0, "", err // use 0 or -1 as dummy status code
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

func randomDelay(min, max int) time.Duration {
	delay := time.Duration(rand.Intn(max-min+1)+min) * time.Second
	time.Sleep(delay)
	return delay
}
