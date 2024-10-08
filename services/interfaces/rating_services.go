package interfaces

import "gotest/models"

type RatingService interface {
	CreateRating(*models.Rating) error
	GetRatingOfMovie(*int) ([]*models.Rating, error)
	GetRatingOfUser(*int) ([]*models.Rating, error)
	GetMovieRatingOfUser(*int, *int) (*models.Rating, error)
	UpdateRating(*models.Rating) error
	DeleteRating(*int, *int) error
	GetAverageRating(*int) (float64, error)
}
