package interfaces

import "gotest/models"

type TopRatedMovieServices interface {
	GetTopRatedMovies() ([]*models.TopRatedMovies, error)
}