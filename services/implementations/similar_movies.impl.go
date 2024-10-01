package implementations

import (
	"context"
	"gotest/models"
	"gotest/services/interfaces"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SimilarMoviesServiceImpl struct {
	similarMoviesCollection *mongo.Collection
	ctx                 context.Context
}

func NewSimilarMoviesService(similarMoviesCollection *mongo.Collection, ctx context.Context) interfaces.SimilarMoviesServices {
	return &SimilarMoviesServiceImpl{
		similarMoviesCollection: similarMoviesCollection,
		ctx:            ctx,
	}
}

func (s * SimilarMoviesServiceImpl) GetSimilarMovies(movieId *int) (*models.SimilarMovies, error) {
	var cast *models.SimilarMovies
	query := bson.D{bson.E{Key: "id", Value: movieId}}
	err := s.similarMoviesCollection.FindOne(s.ctx, query).Decode(&cast)
	return cast, err
}