package purchase

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Store(ctx context.Context, purchase Purchase) error
}

type MongoRepository struct {
	purhcases *mongo.Collection
}

func NewMongoRepo(ctx context.Context, connectionString string) (*MongoRepository, error) {
	return nil, nil
}
