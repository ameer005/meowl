package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ameer005/meowl/internals/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type WebsiteRepo struct {
	collection *mongo.Collection
}

func NewWebsiteRepo(db *mongo.Database) *WebsiteRepo {
	coll := db.Collection("websites")
	createWebsiteIndexes(coll)

	return &WebsiteRepo{
		collection: coll,
	}
}

func (w *WebsiteRepo) GetByURL(ctx context.Context, url string) (*models.Website, error) {
	var website models.Website
	err := w.collection.FindOne(ctx, bson.M{"url": url}).Decode(&website)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &website, nil
}

func (w *WebsiteRepo) AddWebsite(ctx context.Context, data *models.Website) error {
	_, err := w.collection.InsertOne(ctx, data)

	docBytes, _ := bson.Marshal(data)
	log.Printf("Size: %d bytes", len(docBytes))

	return err
}

func createWebsiteIndexes(collection *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "url", Value: 1}}, // Ascending index on "url"
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}
