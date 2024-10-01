package implementations

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gotest/models"
	"gotest/services/interfaces"
)

type CrewServiceImpl struct {
	crewCollection *mongo.Collection
	ctx            context.Context
}

func NewCrewService(crewCollection *mongo.Collection, ctx context.Context) interfaces.CrewService {
	return &CrewServiceImpl{
		crewCollection: crewCollection,
		ctx:            ctx,
	}
}

func (c *CrewServiceImpl) CreateCrew(crew *models.Crew) error {
	_, err := c.crewCollection.InsertOne(c.ctx, crew)
	return err
}

func (c *CrewServiceImpl) GetCrew(movieId *int) (*models.Crew, error) {
	var crew *models.Crew
	query := bson.D{bson.E{Key: "movie_id", Value: movieId}}
	err := c.crewCollection.FindOne(c.ctx, query).Decode(&crew)
	return crew, err
}

func (c *CrewServiceImpl) UpdateCrew(crew *models.Crew) error {
	filter := bson.D{bson.E{Key: "movie_id", Value: crew.MovieId}}
	update := bson.D{
		bson.E{Key: "$set",
			Value: bson.D{
				bson.E{Key: "movie_id", Value: crew.MovieId},
				bson.E{Key: "crew", Value: crew.Crew},
			},
		},
	}
	result, _ := c.crewCollection.UpdateOne(c.ctx, filter, update)
	if result.MatchedCount != 1 {
		return errors.New("no matched document found for update")
	}
	return nil
}

func (c *CrewServiceImpl) DeleteCrew(movieId *int) error {
	filter := bson.D{bson.E{Key: "movie_id", Value: movieId}}
	result, _ := c.crewCollection.DeleteOne(c.ctx, filter)
	if result.DeletedCount != 1 {
		return errors.New("no matched document found for delete")
	}
	return nil
}
