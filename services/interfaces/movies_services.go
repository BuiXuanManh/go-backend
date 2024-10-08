package interfaces

import "gotest/models"

type MovieService interface {
	CreateMovie(*models.Movie) error
	CreateMovies([]*models.Movie) error
	GetMovie(*int) (*models.Movie, error)
	GetMoviesInPage(int, int) ([]*models.Movie, int, error)
	GetPopularMovies(int) ([]*models.Movie, error)
	SearchMovieInPage(*string, *int, *int) ([]*models.Movie, int, error)
	UpdateMovie(*models.Movie) error
	DeleteMovie(*int) error
	FindMovie(*int) (*models.Movie, error)
}
