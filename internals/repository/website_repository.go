package repository

import (
	"context"

	"github.com/ameer005/meowl/internals/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type WebsiteRepo struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewWebsiteRepo(db *mongo.Database) *WebsiteRepo {
	return &WebsiteRepo{
		collection: db.Collection("websites"),
	}
}

func (w *WebsiteRepo) GetByURL(ctx context.Context, url string) (*models.Website, error) {
	var website *models.Website
	err := w.collection.FindOne(ctx, bson.M{"url": url}).Decode(website)

	if err != nil {
		return nil, err
	}

	return website, nil
}

func (w *WebsiteRepo) addWebsite(ctx context.Context, data *models.Website) error {

	return nil
}
