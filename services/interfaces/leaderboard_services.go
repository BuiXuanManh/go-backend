package interfaces

import "gotest/models"

type LeaderboardServices interface {
	GetLeaderboard() ([]*models.Leaderboard, error)
	CreateLeaderboard(*models.Leaderboard) error
	UpdateLeaderboard(*models.Leaderboard, *int, *int) error
	FindLeaderboard(*int, *int) (*models.Leaderboard, error)
}
