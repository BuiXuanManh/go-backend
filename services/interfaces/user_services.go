package interfaces

import "gotest/models"

type UserService interface {
	CreateUser(*models.User) error
	GetUser(*int) (*models.User, error)
	UpdateUser(*models.User) error
	ChangePassword(*int, *string, *string) error
	DeleteUser(*int) error
	GetUserFromUsername(*string) (*models.User, error)
	GetNewUserId() int
}
