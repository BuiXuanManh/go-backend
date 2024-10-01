package interfaces

import "gotest/models"

type UserReviewService interface {
	CreateUserReview(*models.UserReview) error
	GetUserReviewByUserId(*int) ([]*models.UserReview, error)
	GetUserReviewByMovieId(*int) ([]*models.UserReview, error)
	UpdateUserReview(*models.UserReview) error
	DeleteUserReview(*int, *int) error
}