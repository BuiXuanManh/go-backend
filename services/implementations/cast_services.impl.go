package implementations

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gotest/models"
	"gotest/services/interfaces"
)

type CastServiceImpl struct {
	castCollection *mongo.Collection
	ctx            context.Context
}

func NewCastService(castCollection *mongo.Collection, ctx context.Context) interfaces.CastService {
	return &CastServiceImpl{
		castCollection: castCollection,
		ctx:               ctx,
	}
}

func (c *CastServiceImpl) CreateCast(cast *models.Cast) error {
	_, err := c.castCollection.InsertOne(c.ctx, cast)
	return err
}

func (c *CastServiceImpl) GetCast(movieId *int) (*models.Cast, error) {
	var cast *models.Cast
	query := bson.D{bson.E{Key: "movie_id", Value: movieId}}
	err := c.castCollection.FindOne(c.ctx, query).Decode(&cast)
	return cast, err
}

func (c *CastServiceImpl) UpdateCast(cast *models.Cast) error {
	filter := bson.D{bson.E{Key: "movie_id", Value: cast.MovieId}}
	update := bson.D{
		bson.E{Key: "$set",
			Value: bson.D{
				bson.E{Key: "movie_id", Value: cast.MovieId},
				bson.E{Key: "cast", Value: cast.Cast},
			},
		},
	}
	result, _ := c.castCollection.UpdateOne(c.ctx, filter, update)
	if result.MatchedCount != 1 {
		return errors.New("no matched document found for update")
	}
	return nil
}

func (c *CastServiceImpl) DeleteCast(movieId *int) error {
	filter := bson.D{bson.E{Key: "movie_id", Value: movieId}}
	result, _ := c.castCollection.DeleteOne(c.ctx, filter)
	if result.DeletedCount != 1 {
		return errors.New("no matched document found for delete")
	}
	return nil
}
