package implementations

import (
	"context"
	"errors"
	"gotest/models"
	"gotest/services/interfaces"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MovieDiscussionServicesImpl struct {
	movieDiscussionCollection *mongo.Collection
	ctx                       context.Context
}

func NewMovieDiscussionServices(movieDiscussionCollection *mongo.Collection, ctx context.Context) interfaces.MovieDiscussionServices {
	return &MovieDiscussionServicesImpl{
		movieDiscussionCollection: movieDiscussionCollection,
		ctx:                       ctx,
	}
}

func (m *MovieDiscussionServicesImpl) CreateMovieDiscussion(movieDiscussion *models.MovieDiscussion) error {
	_, err := m.movieDiscussionCollection.InsertOne(m.ctx, movieDiscussion)
	return err
}

func (m *MovieDiscussionServicesImpl) GetMovieDiscussion(movieDiscussionId *primitive.ObjectID) (*models.MovieDiscussion, error) {
	var movieDiscussion *models.MovieDiscussion

	query := bson.D{bson.E{Key: "_id", Value: movieDiscussionId}}

	err := m.movieDiscussionCollection.FindOne(m.ctx, query).Decode(&movieDiscussion)

	return movieDiscussion, err
}
func (m *MovieDiscussionServicesImpl) GetMovieDiscussionsByMovieId(movieId *int) ([]*models.MovieDiscussion, error) {
	var movieDiscussions []*models.MovieDiscussion

	query := bson.D{bson.E{Key: "movie_id", Value: movieId}}
	cursor, err := m.movieDiscussionCollection.Find(m.ctx, query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(m.ctx)

	for cursor.Next(m.ctx) {
		var movieDiscussion *models.MovieDiscussion
		if err := cursor.Decode(&movieDiscussion); err != nil {
			return nil, err
		}
		movieDiscussions = append(movieDiscussions, movieDiscussion)
	}
	if len(movieDiscussions) == 0 {
		return nil, errors.New("mongo: no documents in result")
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return movieDiscussions, nil
}

// ------------------------------------------------------------
func (m *MovieDiscussionServicesImpl) GetMovieDiscussionsByUserId(userId *int) ([]*models.MovieDiscussion, error) {
	var movieDiscussions []*models.MovieDiscussion
	query := bson.M{
		"discussion_part.user_id": userId,
	}
	cursor, err := m.movieDiscussionCollection.Find(m.ctx, query)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(m.ctx)

	for cursor.Next(m.ctx) {
		var movieDiscussion *models.MovieDiscussion
		if err := cursor.Decode(&movieDiscussion); err != nil {
			return nil, err
		}
		movieDiscussions = append(movieDiscussions, movieDiscussion)
	}
	if len(movieDiscussions) == 0 {
		return nil, errors.New("mongo: no documents in result")
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return movieDiscussions, nil
}

//----------------------------------------------------------

func (m *MovieDiscussionServicesImpl) UpdateMovieDiscussion(movieDiscussion *models.MovieDiscussion) error {
	filter := bson.D{bson.E{Key: "_id", Value: movieDiscussion.ID}}
	update := bson.D{
		bson.E{Key: "$set",
			Value: bson.D{
				bson.E{Key: "subject", Value: movieDiscussion.Subject},
				bson.E{Key: "status", Value: movieDiscussion.Status},
				bson.E{Key: "discussion_part", Value: movieDiscussion.DiscussionPart},
			},
		},
	}
	result, _ := m.movieDiscussionCollection.UpdateOne(m.ctx, filter, update)
	if result.MatchedCount != 1 {
		return errors.New("no matched document found for update")
	}
	return nil
}

func (m *MovieDiscussionServicesImpl) DeleteMovieDiscussion(movieDiscussionId *primitive.ObjectID) error {
	filter := bson.D{bson.E{Key: "_id", Value: movieDiscussionId}}
	result, _ := m.movieDiscussionCollection.DeleteOne(m.ctx, filter)
	if result.DeletedCount != 1 {
		return errors.New("no matched document found for delete")
	}
	return nil
}

func (m *MovieDiscussionServicesImpl) CreateMovieDiscussionPart(discussionPart *models.DiscussionPart, discussionId *primitive.ObjectID) error {
	// Fetch the MovieDiscussion document
	var movieDiscussion models.MovieDiscussion
	filter := bson.M{"_id": discussionId}
	err := m.movieDiscussionCollection.FindOne(m.ctx, filter).Decode(&movieDiscussion)
	if err != nil {
		return err
	}

	// Calculate the next part_id based on the count of existing DiscussionPart elements
	nextPartID := len(movieDiscussion.DiscussionPart)
	discussionPart.PartId = nextPartID

	// Update the MovieDiscussion document with the new DiscussionPart
	update := bson.M{
		"$push": bson.M{"discussion_part": discussionPart},
	}
	_, err = m.movieDiscussionCollection.UpdateOne(m.ctx, filter, update)
	return err
}

func (m *MovieDiscussionServicesImpl) GetMovieDiscussionInPage(pageNumber int, discussionPerPage int) ([]*models.MovieDiscussion, int, error) {
	var movieDiscussionInPage []*models.MovieDiscussion
	options := options.Find().
		SetSkip(int64(pageNumber - 1) * int64(pageNumber)).
		SetLimit(int64(discussionPerPage))

	cursor, err := m.movieDiscussionCollection.Find(m.ctx, bson.D{{}}, options)

	if err != nil {
		return nil, 0, err
	}
	for cursor.Next(m.ctx) {
		var movieDiscussion models.MovieDiscussion
		err := cursor.Decode(&movieDiscussion)
		if err != nil {
			return nil, 0, err
		}
		movieDiscussionInPage = append(movieDiscussionInPage, &movieDiscussion)
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}
	cursor.Close(m.ctx)
	if len(movieDiscussionInPage) == 0 {
		return nil, 0, errors.New("documents not found")
	}

	// Fetch total movies
	totalMovies, err := m.movieDiscussionCollection.CountDocuments(m.ctx, bson.D{{}})
	if err != nil {
		return nil, 0, err
	}

	return movieDiscussionInPage, int(totalMovies), nil
}

func (m *MovieDiscussionServicesImpl) UpdateMovieDiscussionPart(discussionId *primitive.ObjectID, partId *int, updatedPart *models.DiscussionPart) error {
	filter := bson.M{
		"_id":                     discussionId,
		"discussion_part.part_id": partId,
	}

	update := bson.M{
		"$set": bson.M{
			"discussion_part.$": updatedPart,
		},
	}

	result, err := m.movieDiscussionCollection.UpdateOne(m.ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount != 1 {
		return errors.New("no matched document found for update")
	}

	return nil
}

func (m *MovieDiscussionServicesImpl) DeleteMovieDiscussionPart(discussionId *primitive.ObjectID, partId *int) error {
	filter := bson.M{
		"_id":                     discussionId,
		"discussion_part.part_id": partId,
	}

	update := bson.M{
		"$pull": bson.M{
			"discussion_part": bson.M{"part_id": partId},
		},
	}

	result, err := m.movieDiscussionCollection.UpdateOne(m.ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount != 1 {
		return errors.New("no matched document found for delete")
	}

	return nil
}
