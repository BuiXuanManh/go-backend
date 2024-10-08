package implementations

import (
	"context"
	"errors"
	"gotest/models"
	"gotest/services/interfaces"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LeaderboardServiceImpl struct {
	leaderboardCollection *mongo.Collection
	ctx                   context.Context
}

// FindLeaderboard implements interfaces.LeaderboardServices.
func (l *LeaderboardServiceImpl) FindLeaderboard(movie_id *int, user_id *int) (*models.Leaderboard, error) {
	var leaderboard models.Leaderboard
	filter := bson.M{"movie_id": movie_id, "user_id": user_id}
	err := l.leaderboardCollection.FindOne(l.ctx, filter).Decode(&leaderboard)
	if err == mongo.ErrNoDocuments {
		// Nếu không tìm thấy tài liệu, trả về nil mà không dừng chương trình
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &leaderboard, nil
}

// CreateLeaderboard implements interfaces.LeaderboardServices.
func (l *LeaderboardServiceImpl) CreateLeaderboard(new *models.Leaderboard) error {
	_, err := l.leaderboardCollection.InsertOne(l.ctx, new)
	if err != nil {
		return err
	}
	return nil
}

// UpdateLeaderboard implements interfaces.LeaderboardServices.
func (l *LeaderboardServiceImpl) UpdateLeaderboard(lboard *models.Leaderboard, movie_id *int, user_id *int) error {
	filter := bson.M{"movie_id": movie_id, "user_id": user_id}
	update := bson.M{"$set": bson.M{"movies_rated": lboard.MoviesRated}}
	_, err := l.leaderboardCollection.UpdateOne(l.ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func NewLeaderboardService(leaderboardCollection *mongo.Collection, ctx context.Context) interfaces.LeaderboardServices {
	return &LeaderboardServiceImpl{
		leaderboardCollection: leaderboardCollection,
		ctx:                   ctx,
	}
}
func (l *LeaderboardServiceImpl) GetLeaderboard() ([]*models.Leaderboard, error) {
	var leaderboardUsers []*models.Leaderboard
	cursor, err := l.leaderboardCollection.Find(l.ctx, bson.D{{}})

	if err != nil {
		return nil, err
	}
	for cursor.Next(l.ctx) {
		var leaderboardUser models.Leaderboard
		err := cursor.Decode(&leaderboardUser)
		if err != nil {
			return nil, err
		}
		leaderboardUsers = append(leaderboardUsers, &leaderboardUser)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	cursor.Close(l.ctx)

	if len(leaderboardUsers) == 0 {
		return nil, errors.New("documents not found")
	}

	return leaderboardUsers, nil
}
