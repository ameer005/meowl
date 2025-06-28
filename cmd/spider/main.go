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

	postgresClient, err := storage.NewPostgresConnection(os.Getenv("POSTGRES_URL"))
	if err != nil {
		os.Exit(1)
	}
	defer postgresClient.Close()

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

	crawler := crawler.New(input, logger.Logger, postgresClient)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for i := 0; i < 10; i++ {
		crawler.Wg.Add(1)
		go crawler.Start(ctx)
	}
	crawler.Wg.Wait()
}
