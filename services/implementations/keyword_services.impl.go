package implementations

import (
	"context"
	"gotest/models"
	"gotest/services/interfaces"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type KeywordServiceImpl struct {
	keywordCollection *mongo.Collection
	ctx context.Context
}

func NewKeywordService(keywordCollection *mongo.Collection, ctx context.Context) interfaces.KeywordService{
	return &KeywordServiceImpl{
		keywordCollection: keywordCollection,
		ctx: ctx,
	}
}

func (k *KeywordServiceImpl) CreateKeyword(keyword *models.Keyword) error{
	_, err := k.keywordCollection.InsertOne(k.ctx, keyword)
	return err
}

func (k *KeywordServiceImpl) GetKeyword(movieId *int) (*models.Keyword, error){
	var keyword *models.Keyword
	query := bson.D{bson.E{Key: "movie_id", Value: movieId}}
	err := k.keywordCollection.FindOne(k.ctx, query).Decode(&keyword)
	return keyword, err
}

func (k *KeywordServiceImpl) UpdateKeyword(keyword *models.Keyword) error{
	filter := bson.D{bson.E{Key: "movie_id", Value: keyword.MovieId}}
	update := bson.D{
		bson.E{Key: "$set", 
		Value: bson.D{
			bson.E{Key: "movie_id", Value: keyword.MovieId},
			bson.E{Key: "keyword_list", Value: keyword.KeywordList},
			},
		},
	}
	result, _ := k.keywordCollection.UpdateOne(k.ctx, filter, update)
	if result.MatchedCount != 1 {
		return errors.New("no matched document found for update")
	}
	return nil
}
func (k *KeywordServiceImpl) DeleteKeyword(movieId *int) error{
	filter := bson.D{bson.E{Key: "movie_id", Value: movieId}}
	result, _ := k.keywordCollection.DeleteOne(k.ctx, filter)
	if result.DeletedCount != 1{
		return errors.New("no matched document found for delete")
	}
	return nil
}