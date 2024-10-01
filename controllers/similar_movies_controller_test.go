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

type SimilarMoviesTestSuite struct {
	suite.Suite
	r                       *gin.Engine
	SimilarMoviesCollection *mongo.Collection
	SimilarMoviesService    interfaces.SimilarMoviesServices
	SimilarMoviesController SimilarMoviesController
	setUpDone               bool
}

func (suite *SimilarMoviesTestSuite) SetupTest() {
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

	suite.SimilarMoviesCollection = mongoClient.Database(os.Getenv("DB_NAME")).Collection("similar_movies")
	suite.SimilarMoviesService = implementations.NewSimilarMoviesService(suite.SimilarMoviesCollection, ctx)
	suite.SimilarMoviesController = NewSimilarMoviesController(suite.SimilarMoviesService)

	suite.setUpDone = true
}

func (suite *SimilarMoviesTestSuite) setupRouter() *gin.Engine {
	r := gin.New()
	r.GET("/similarMovies/get/:id", suite.SimilarMoviesController.GetSimilarMovies)
	return r
}

func (suite *SimilarMoviesTestSuite) TestGetSimilarMoviesSuccessfully() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/similarMovies/get/1024", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	assert.Equal(suite.T(), `{"id":1024,"similar_movies":[1024,2029,2860,2267,1273,253,527,100017,4011,42002]}`, w.Body.String())
}

func (suite *SimilarMoviesTestSuite) TestGetNonExistingSimilarMovie() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/similarMovies/get/1", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"mongo: no documents in result"}`, w.Body.String())
}

func (suite *SimilarMoviesTestSuite) TestGetInvalidMovieId() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/similarMovies/get/fsaf", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"strconv.Atoi: parsing \"fsaf\": invalid syntax"}`, w.Body.String())
}

func (suite *SimilarMoviesTestSuite) TestRegisterRouteSuccessfully() {
	// Create a new gin router
	router := gin.New()

	// Create an instance of your UserReviewController
	cc := suite.SimilarMoviesController // replace UserReviewController with your actual type

	// Register the UserReview routes
	cc.RegisterSimilarMoviesRoute(router.Group("/v2"))

	// Perform a GET request to "/v2/UserReview/get/1"
	w := performSimilarMoviesRequest(router, "GET", "/v2/similarMovie/get/1024")

	// Assert the response status code
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// You can add more assertions for other routes if needed
}

func performSimilarMoviesRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestSimilarMoviesTestSuite(t *testing.T) {
	suite.Run(t, new(SimilarMoviesTestSuite))
}
