package interfaces

import "gotest/models"

type CastService interface {
	CreateCast(*models.Cast) error
	GetCast(*int) (*models.Cast, error)
	UpdateCast(*models.Cast) error
	DeleteCast(*int) error
}
