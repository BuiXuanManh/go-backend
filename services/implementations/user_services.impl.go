package implementations

import (
	"context"
	"errors"
	"gotest/helper"
	"gotest/models"
	"gotest/services/interfaces"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	userCollection *mongo.Collection
	ctx            context.Context
}

func NewUserService(userCollection *mongo.Collection, ctx context.Context) interfaces.UserService {
	return &UserServiceImpl{
		userCollection: userCollection,
		ctx:            ctx,
	}
}

func (u *UserServiceImpl) CreateUser(user *models.User) error {
	_, err := u.userCollection.InsertOne(u.ctx, user)
	return err
}

func (u *UserServiceImpl) GetUser(userId *int) (*models.User, error) {
	var user *models.User
	query := bson.D{bson.E{Key: "id", Value: userId}}
	err := u.userCollection.FindOne(u.ctx, query).Decode(&user)
	return user, err
}

func (u *UserServiceImpl) UpdateUser(user *models.User) error {
	filter := bson.D{bson.E{Key: "id", Value: user.UserId}}
	update := bson.D{
		bson.E{Key: "$set",
			Value: bson.D{
				bson.E{Key: "username", Value: user.Username},
				bson.E{Key: "password_hash", Value: user.PasswordHash},
				bson.E{Key: "favorite_list", Value: user.FavoriteList},
				bson.E{Key: "recommendation_list", Value: user.RecommendationList},
				bson.E{Key: "watch_list", Value: user.WatchList},
				bson.E{Key: "picture_profile", Value: user.PictureProfile},
			},
		},
	}
	result, _ := u.userCollection.UpdateOne(u.ctx, filter, update)
	if result.MatchedCount != 1 {
		return errors.New("no matched document found for update")
	}
	return nil
}

func (u *UserServiceImpl) ChangePassword(userId *int, OldPassword *string, NewPassword *string) error {
	// Retrieve the user by ID
	filter := bson.D{bson.E{Key: "id", Value: userId}}
	var existingUser *models.User
	err := u.userCollection.FindOne(u.ctx, filter).Decode(&existingUser)
	if err != nil {
		return err
	}
	// Compare password with password hash
	if !helper.CheckPassword(existingUser.PasswordHash, *OldPassword) {
		return errors.New("wrong password")
	}
	hashedPassword, err := helper.HashPassword(*NewPassword)
	if err != nil {
		return err
	}
	existingUser.PasswordHash = string(hashedPassword)

	update := bson.D{
		bson.E{Key: "$set",

			Value: bson.D{
				bson.E{Key: "password_hash", Value: existingUser.PasswordHash},
			},
		},
	}

	result, _ := u.userCollection.UpdateOne(u.ctx, filter, update)
	if result.MatchedCount != 1 {
		return errors.New("no matched document found for update")
	}
	return nil
}

func (u *UserServiceImpl) DeleteUser(userId *int) error {
	filter := bson.D{bson.E{Key: "id", Value: userId}}
	result, _ := u.userCollection.DeleteOne(u.ctx, filter)
	if result.DeletedCount != 1 {
		return errors.New("no matched document found for delete")
	}
	return nil
}

func (u *UserServiceImpl) GetUserFromUsername(username *string) (*models.User, error) {
	var user *models.User
	query := bson.D{bson.E{Key: "username", Value: username}}
	err := u.userCollection.FindOne(u.ctx, query).Decode(&user)
	return user, err
}

func (u *UserServiceImpl) GetNewUserId() int {
	// Fetch total movies
	totalUsers, err := u.userCollection.CountDocuments(u.ctx, bson.D{{}})
	if err != nil {
		return 0
	}

	return int(totalUsers)
}
