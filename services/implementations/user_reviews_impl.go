package implementations

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gotest/models"
	"gotest/services/interfaces"
)

type UserReviewServiceImpl struct {
	userReviewCollection *mongo.Collection
	ctx            context.Context
}

func NewUserReviewService(userReviewCollection *mongo.Collection, ctx context.Context) interfaces.UserReviewService {
	return &UserReviewServiceImpl{
		userReviewCollection: userReviewCollection,
		ctx:               ctx,
	}
}

func (u *UserReviewServiceImpl) CreateUserReview(userReview *models.UserReview) error {
	_, err := u.userReviewCollection.InsertOne(u.ctx, userReview)
	return err
}

func (u *UserReviewServiceImpl) GetUserReviewByUserId(userId *int) ([]*models.UserReview, error) {
	var userReviews []*models.UserReview
	query := bson.M{
		"user_id": userId,
	}
	cursor, err := u.userReviewCollection.Find(u.ctx, query)
	
	if err != nil {
		return nil, err
	}
	defer cursor.Close(u.ctx)

	for cursor.Next(u.ctx) {
		var userReview *models.UserReview
		if err := cursor.Decode(&userReview); err != nil {
			return nil, err
		}
		userReviews = append(userReviews, userReview)
	}
	if len(userReviews) == 0 {
		return nil, errors.New("mongo: no documents in result")
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return userReviews, nil
}

func (u *UserReviewServiceImpl) GetUserReviewByMovieId(movieId *int) ([]*models.UserReview, error) {
	var userReviews []*models.UserReview
	query := bson.M{
		"movie_id": movieId,
	}
	cursor, err := u.userReviewCollection.Find(u.ctx, query)
	
	if err != nil {
		return nil, err
	}
	defer cursor.Close(u.ctx)

	for cursor.Next(u.ctx) {
		var userReview *models.UserReview
		if err := cursor.Decode(&userReview); err != nil {
			return nil, err
		}
		userReviews = append(userReviews, userReview)
	}
	if len(userReviews) == 0 {
		return nil, errors.New("mongo: no documents in result")
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return userReviews, nil
}

func (u *UserReviewServiceImpl) UpdateUserReview(userReview *models.UserReview) error {
	filter := bson.M{
		"user_id": userReview.UserID,
		"movie_id": userReview.MovieID,
	}
	update := bson.D{
		bson.E{Key: "$set",
			Value: bson.D{
				bson.E{Key: "comment", Value: userReview.Comment},
			},
		},
	}

	result, _ := u.userReviewCollection.UpdateOne(u.ctx, filter, update)
	if result.MatchedCount != 1 {
		return errors.New("no matched document found for update")
	}
	return nil
}

func (u *UserReviewServiceImpl) DeleteUserReview(userId, movieId *int) error {
	filter := bson.M{
		"user_id": userId,
		"movie_id": movieId,
	}
	result, _ := u.userReviewCollection.DeleteOne(u.ctx, filter)
	if result.DeletedCount != 1 {
		return errors.New("no matched document found for delete")
	}
	return nil
}
