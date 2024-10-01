package controllers

import (
	"context"
	"log"
	"gotest/services/implementations"
	"gotest/services/interfaces"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	// "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	// "github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type LeaderboardTestSuite struct{
	suite.Suite
	r *gin.Engine
	LeaderboardCollection *mongo.Collection
	LeaderboardService interfaces.LeaderboardServices
	LeaderboardController LeaderboardController
	setUpDone bool
}

func (suite *LeaderboardTestSuite) SetupTest() {
	if suite.setUpDone {
		return
	}

	ctx := context.TODO()
 
	// err := godotenv.Load("../.env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	// Retrieve the connection string from the environment
	connectionString := os.Getenv("DB_CONNECTION_STRING")

	mongoConn := options.
	Client().
	ApplyURI(connectionString)

	mongoClient, err := mongo.Connect(ctx, mongoConn)
	if err != nil {
		log.Fatal(err)
	}
	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	suite.r = suite.setupRouter()

    suite.LeaderboardCollection = mongoClient.Database(os.Getenv("DB_NAME")).Collection("leaderboard")
	suite.LeaderboardService = implementations.NewLeaderboardService(suite.LeaderboardCollection, ctx)
	suite.LeaderboardController = NewLeaderboardController(suite.LeaderboardService)

	suite.setUpDone = true
}

func (suite *LeaderboardTestSuite) setupRouter() *gin.Engine {
	r := gin.New()
	r.GET("/leaderboard/get", suite.LeaderboardController.GetLeaderboard)
	return r
}

func (suite *LeaderboardTestSuite) TestGetLeaderboardSuccessfully() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/leaderboard/get", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusOK, w.Code)	
}

func (suite *LeaderboardTestSuite) TestRegisterRouteSuccessfully() {
	// Create a new gin router
	router := gin.New()

	// Create an instance of your UserReviewController
	cc := suite.LeaderboardController // replace UserReviewController with your actual type

	// Register the UserReview routes
	cc.RegisterLeaderboardRoute(router.Group("/v2"))

	// Perform a GET request to "/v2/UserReview/get/1"
	w := performUserReviewRequest(router, "GET", "/v2/leaderboard/get")

	// Assert the response status code
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// You can add more assertions for other routes if needed
}

func TestLeaderboardTestSuite (t *testing.T){
	suite.Run(t, new(LeaderboardTestSuite))
}