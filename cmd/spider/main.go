package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ameer005/meowl/internals/crawler"
	"github.com/ameer005/meowl/pkg/logger"
)

func main() {
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

	crawl := crawler.New(input, logger.Logger)
	crawl.Start()
}
