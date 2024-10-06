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

func NewLeaderboardService(leaderboardCollection *mongo.Collection, ctx context.Context) interfaces.LeaderboardServices {
	return &LeaderboardServiceImpl{
		leaderboardCollection: leaderboardCollection,
		ctx:                   ctx,
	}
}
func (l *LeaderboardServiceImpl) updateLeaderboard() (*[]models.Leaderboard, error) {
	//thống kê điểm số rate của của movie

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
