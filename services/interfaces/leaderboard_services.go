package interfaces

import "gotest/models"

type LeaderboardServices interface {
	GetLeaderboard() ([]*models.Leaderboard, error)
	updateLeaderboard() (*[]models.Leaderboard, error)
}
