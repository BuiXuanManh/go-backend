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
	ctx                      context.Context
}

// UpdateTopRatedMovies implements interfaces.TopRatedMovieServices.
func (l *TopRatedMoviesImpl) UpdateTopRatedMovies(top *models.TopRatedMovies) error {
	filter := bson.M{"movie_id": top.MovieId}
	update := bson.M{"$set": &top}
	_, err := l.topRatedMoviesCollection.UpdateOne(l.ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// GetAverageRating implements interfaces.TopRatedMovieServices.

// FindTopRatedMovies implements interfaces.TopRatedMovieServices.
func (l *TopRatedMoviesImpl) FindTopRatedMovies(movies_id *int) (*models.TopRatedMovies, error) {
	var topRatedMovie models.TopRatedMovies
	filter := bson.M{"movie_id": movies_id}
	err := l.topRatedMoviesCollection.FindOne(l.ctx, filter).Decode(&topRatedMovie)
	if err == mongo.ErrNoDocuments {
		// Nếu không tìm thấy tài liệu, trả về nil mà không dừng chương trình
		return nil, nil
	}

	// Xử lý các lỗi khác nếu có
	if err != nil {
		return nil, err
	}

	return &topRatedMovie, nil
}

// CreateTopRatedMovies implements interfaces.TopRatedMovieServices.
func (l *TopRatedMoviesImpl) CreateTopRatedMovies(top *models.TopRatedMovies) error {
	_, err := l.topRatedMoviesCollection.InsertOne(l.ctx, top)
	if err != nil {
		return err
	}
	return nil
}

func NewTopRatedMoviesService(topRatedMoviesCollection *mongo.Collection, ctx context.Context) interfaces.TopRatedMovieServices {
	return &TopRatedMoviesImpl{
		topRatedMoviesCollection: topRatedMoviesCollection,
		ctx:                      ctx,
	}
}

func (l *TopRatedMoviesImpl) GetTopRatedMovies() ([]*models.TopRatedMovies, error) {
	var topRatedMovies []*models.TopRatedMovies
	cursor, err := l.topRatedMoviesCollection.Find(l.ctx, bson.D{{}})

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
