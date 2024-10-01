package interfaces

import "gotest/models"

type SimilarMoviesServices interface {
	GetSimilarMovies(*int) (*models.SimilarMovies, error)
}