package implementations

import (
	"context"
	"errors"
	"gotest/models"
	"gotest/services/interfaces"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RatingServiceImpl struct {
	ratingCollection *mongo.Collection
	ctx              context.Context
}

func NewRatingService(ratingCollection *mongo.Collection, ctx context.Context) interfaces.RatingService {
	return &RatingServiceImpl{
		ratingCollection: ratingCollection,
		ctx:              ctx,
	}
}

func (r *RatingServiceImpl) CreateRating(rating *models.Rating) error {
	_, err := r.ratingCollection.InsertOne(r.ctx, rating)
	return err
}

func (r *RatingServiceImpl) UpdateRating(rating *models.Rating) error {
	filter := bson.D{
		bson.E{Key: "movie_id", Value: rating.MovieId}, bson.E{Key: "user_id", Value: rating.UserId},
	}
	update := bson.D{
		bson.E{Key: "$set",
			Value: bson.D{
				bson.E{Key: "rating", Value: rating.Rating},
				bson.E{Key: "timestamp", Value: rating.Timestamp},
			},
		},
	}
	result, _ := r.ratingCollection.UpdateOne(r.ctx, filter, update)
	if result.MatchedCount != 1 {
		return errors.New("no matched document found for update")
	}
	return nil
}

func (r *RatingServiceImpl) DeleteRating(movieId *int, userId *int) error {
	filter := bson.D{
		bson.E{Key: "movie_id", Value: movieId}, bson.E{Key: "user_id", Value: userId},
	}
	result, _ := r.ratingCollection.DeleteOne(r.ctx, filter)
	if result.DeletedCount != 1 {
		return errors.New("no matched document found for delete")
	}
	return nil
}

func (r *RatingServiceImpl) GetRatingOfMovie(movieId *int) ([]*models.Rating, error) {
	var ratings []*models.Rating
	filter := bson.D{
		bson.E{Key: "movie_id", Value: movieId},
	}
	cursor, err := r.ratingCollection.Find(r.ctx, filter)

	if err != nil {
		return nil, err
	}
	for cursor.Next(r.ctx) {
		var rating models.Rating
		err := cursor.Decode(&rating)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, &rating)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	cursor.Close(r.ctx)

	if len(ratings) == 0 {
		return nil, errors.New("documents not found")
	}

	return ratings, nil
}

func (r *RatingServiceImpl) GetRatingOfUser(userId *int) ([]*models.Rating, error) {
	var ratings []*models.Rating
	filter := bson.D{
		bson.E{Key: "user_id", Value: userId},
	}
	cursor, err := r.ratingCollection.Find(r.ctx, filter)

	if err != nil {
		return nil, err
	}
	for cursor.Next(r.ctx) {
		var rating models.Rating
		err := cursor.Decode(&rating)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, &rating)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	cursor.Close(r.ctx)

	if len(ratings) == 0 {
		return nil, errors.New("documents not found")
	}

	return ratings, nil
}

func (r *RatingServiceImpl) GetMovieRatingOfUser(movieId *int, userId *int) (*models.Rating, error) {
	var rating models.Rating

	filter := bson.D{
		bson.E{Key: "movie_id", Value: movieId},
		bson.E{Key: "user_id", Value: userId},
	}

	err := r.ratingCollection.FindOne(r.ctx, filter).Decode(&rating)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("rating not found")
		}
		return nil, err
	}

	return &rating, nil
}
