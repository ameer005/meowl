package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ameer005/meowl/internals/crawler"
	"github.com/ameer005/meowl/internals/storage"
	"github.com/ameer005/meowl/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	mongoClient, err := storage.NewMongoConnection(os.Getenv("MONGO_URL"))

	if err != nil {
		os.Exit(1)
	}

	logger.Init()
	fs := flag.NewFlagSet("spider", flag.ExitOnError)

	var inputRaw string
	fs.StringVar(&inputRaw, "crawl", "", "Comma-separated list of input URLs")

	fs.Parse(os.Args[1:])

	if inputRaw == "" {
		fmt.Println("Usage: go run main.go -crawl=url1,url2,...")
		return
	}

	input := strings.Split(inputRaw, ",")

	crawler := crawler.New(input, logger.Logger, mongoClient)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for i := 0; i < 10; i++ {
		crawler.Wg.Add(1)
		go crawler.Start(ctx)
	}
	crawler.Wg.Wait()
}
