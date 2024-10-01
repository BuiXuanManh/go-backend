package interfaces

import (
	"gotest/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MovieDiscussionServices interface {
	CreateMovieDiscussion(*models.MovieDiscussion) error
	GetMovieDiscussion(*primitive.ObjectID) (*models.MovieDiscussion, error)
	GetMovieDiscussionsByMovieId(*int) ([]*models.MovieDiscussion, error)
	GetMovieDiscussionsByUserId(*int) ([]*models.MovieDiscussion, error)
	UpdateMovieDiscussion(*models.MovieDiscussion) error
	DeleteMovieDiscussion(*primitive.ObjectID) error
	CreateMovieDiscussionPart(*models.DiscussionPart, *primitive.ObjectID) error
	GetMovieDiscussionInPage(int, int) ([]*models.MovieDiscussion, int, error)
	UpdateMovieDiscussionPart(*primitive.ObjectID, *int, *models.DiscussionPart) error
	DeleteMovieDiscussionPart(*primitive.ObjectID, *int) error
}
