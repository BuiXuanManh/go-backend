package interfaces

import "gotest/models"

type KeywordService interface {
	CreateKeyword(*models.Keyword) error
	GetKeyword(*int) (*models.Keyword, error)
	UpdateKeyword(*models.Keyword) error
	DeleteKeyword(*int) error
}