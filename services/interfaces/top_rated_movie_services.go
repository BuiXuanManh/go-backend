package interfaces

import "gotest/models"

type TopRatedMovieServices interface {
	GetTopRatedMovies() ([]*models.TopRatedMovies, error)
	CreateTopRatedMovies(*models.TopRatedMovies) error
	UpdateTopRatedMovies(*models.TopRatedMovies) error
	FindTopRatedMovies(*int) (*models.TopRatedMovies, error)
	// GetAverageRating(*int) (models.TopRatedMovies, error)
}
