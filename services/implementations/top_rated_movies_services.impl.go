package implementations

import (
	"context"
	"errors"
	"gotest/models"
	"gotest/services/interfaces"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TopRatedMoviesImpl struct {
	topRatedMoviesCollection *mongo.Collection
	ctx                 context.Context
}

func NewTopRatedMoviesService(topRatedMoviesCollection *mongo.Collection, ctx context.Context) interfaces.TopRatedMovieServices {
	return &TopRatedMoviesImpl{
		topRatedMoviesCollection: topRatedMoviesCollection,
		ctx:            ctx,
	}
}

func (l *TopRatedMoviesImpl) GetTopRatedMovies() ([]*models.TopRatedMovies, error){
	var topRatedMovies []*models.TopRatedMovies
	cursor, err :=  l.topRatedMoviesCollection.Find(l.ctx, bson.D{{}})
	
	if err != nil {
		return nil, err
	}
	for cursor.Next(l.ctx) {
		var topRatedMovie models.TopRatedMovies
		err := cursor.Decode(&topRatedMovie)
		if err != nil {
			return nil, err
		}
		topRatedMovies = append(topRatedMovies, &topRatedMovie)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	cursor.Close(l.ctx)

	if len(topRatedMovies) == 0 {
		return nil, errors.New("documents not found")
	}

	return topRatedMovies, nil
}