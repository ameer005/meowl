package main

import (
	"os"

	"github.com/ameer005/meowl/internals/repository"
	"github.com/ameer005/meowl/internals/storage"
	"github.com/ameer005/meowl/internals/tokenizer"
	"github.com/ameer005/meowl/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	logger.Init()

	postgresClient, err := storage.NewPostgresConnection(os.Getenv("POSTGRES_URL"))
	if err != nil {
		os.Exit(1)
	}
	defer postgresClient.Close()

	lastSeenId := 319
	batch := 1

	spiderRepo := repository.NewPostgresSpiderWebsiteRepo(postgresClient)

	spiderData, err := spiderRepo.GetWebsitesBatch(lastSeenId, batch)
	if err != nil {
		logger.Logger.Error("Failed to fetch spider data")
		os.Exit(1)
	}

	tokenizer.Tokenize(spiderData[0].Content)

}
