package controllers

import (
	"context"
	"log"
	"gotest/services/implementations"
	"gotest/services/interfaces"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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

type RatingTestSuite struct {
	suite.Suite
	r                *gin.Engine
	ratingCollection *mongo.Collection
	ratingService    interfaces.RatingService
	ratingController RatingController
	setUpDone        bool
}

func (suite *RatingTestSuite) SetupTest() {
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
	suite.r = suite.setupRatingRouter()

	suite.ratingCollection = mongoClient.Database(os.Getenv("DB_NAME")).Collection("ratings")
	suite.ratingService = implementations.NewRatingService(suite.ratingCollection, ctx)
	suite.ratingController = NewRatingController(suite.ratingService)

	suite.setUpDone = true
}

func (suite *RatingTestSuite) TestCreateRatingSuccessfully() {
	body := strings.NewReader(`{
		"user_id": 1000,
		"movie_id": 1,
		"rating": 5.0,
		"timestamp": 1260759135000
	  }`)

	req, _ := http.NewRequest("POST", "/rating/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	assert.Equal(suite.T(), `{"message":"Successful"}`, w.Body.String())
}

func (suite *RatingTestSuite) TestCreateExistingRating() {
	body := strings.NewReader(`{
		"user_id": 1,
		"movie_id": 1371,
		"rating": 2.5,
		"timestamp": 1260759135000
	  }`)

	req, _ := http.NewRequest("POST", "/rating/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	assert.Equal(suite.T(), `{"message":"write exception: write errors: [E11000 duplicate key error collection: moviesDB.ratings index: movie_id_1_user_id_1 dup key: { movie_id: 1371, user_id: 1 }]"}`, w.Body.String())
}

func (suite *RatingTestSuite) TestCreateWrongBindedJSONRating() {
	body := strings.NewReader(`{UnbindedJSON}`)

	req, _ := http.NewRequest("POST", "/rating/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"invalid character 'U' looking for beginning of object key string"}`, w.Body.String())
}

func (suite *RatingTestSuite) TestCreateWrongBodyStructureInputRating() {
	body := strings.NewReader(`{"test": true}`)

	req, _ := http.NewRequest("POST", "/rating/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"wrong input structure"}`, w.Body.String())
}

// func (suite *RatingTestSuite) TestGetAnExistingRating() {
// 	// Create a fake HTTP request
// 	req, _ := http.NewRequest("GET", "/keyword/get/1", nil)

// 	// Create a response recorder to record the response
// 	w := httptest.NewRecorder()

// 	// Perform the request
// 	suite.r.ServeHTTP(w, req)

// 	// Check the response status code
// 	assert.Equal(suite.T(), http.StatusOK, w.Code)

// 	// Check the response body
// 	assert.JSONEq(suite.T(), `{
// 		"movie_id": 1,
// 		"keyword_list": [
// 		  "robbery",
// 		  "detective",
// 		  "bank",
// 		  "obsession"]}`, w.Body.String())
// }

func (suite *RatingTestSuite) TestGetANonExistingUserId() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/rating/get/user/10000", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"documents not found"}`, w.Body.String())
}

func (suite *RatingTestSuite) TestGetInvalidUserIdRating() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/rating/get/user/fsaf", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"Invalid movie id"}`, w.Body.String())
}

func (suite *RatingTestSuite) TestGetANonExistingMovieId() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/rating/get/movie/100000", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"documents not found"}`, w.Body.String())
}

func (suite *RatingTestSuite) TestGetInvalidMovieIdRating() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/rating/get/movie/fsaf", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"Invalid movie id"}`, w.Body.String())
}

func (suite *RatingTestSuite) TestUpdateUnbindJSONRating() {
	body := strings.NewReader(`{UnbindedJSON}`)

	req, _ := http.NewRequest("PATCH", "/rating/update", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"invalid character 'U' looking for beginning of object key string"}`, w.Body.String())
}

func (suite *RatingTestSuite) TestUpdateWrongStructureRating() {
	body := strings.NewReader(`{"test": true}`)

	req, _ := http.NewRequest("PATCH", "/rating/update", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"wrong input structure"}`, w.Body.String())
}

func (suite *RatingTestSuite) TestUpdateNonExistentRating() {
	body := strings.NewReader(`{
		"user_id": 20000,
		"movie_id": 1,
		"rating": 5.0,
		"timestamp": 1260759135000
	  }`)
	// Create a fake HTTP request
	req, _ := http.NewRequest("PATCH", "/rating/update", body)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"no matched document found for update"}`, w.Body.String())
}

func (suite *RatingTestSuite) TestUpdateSuccessfulRating() {
	body := strings.NewReader(`{
		"user_id": 1000,
		"movie_id": 1,
		"rating": 4.5,
		"timestamp": 1260759135000
	  }`)

	req, _ := http.NewRequest("PATCH", "/rating/update", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	assert.Equal(suite.T(), `{"message":"Successful"}`, w.Body.String())
}

func (suite *RatingTestSuite) TestDeleteInvalidMovieIdRating() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("DELETE", "/rating/delete/1/fsaf", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"Invalid movie id"}`, w.Body.String())
}

func (suite *RatingTestSuite) TestDeleteInvalidUserIdRating() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("DELETE", "/rating/delete/fsaf/1", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"Invalid user id"}`, w.Body.String())
}

func (suite *RatingTestSuite) TestDeleteNonExistentUserId() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("DELETE", "/rating/delete/100000/1", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"no matched document found for delete"}`, w.Body.String())
}

func (suite *RatingTestSuite) TestRegisterRatingRouteSuccessfully() {
	// Create a new gin router
	router := gin.New()

	// Create an instance of your CastController
	cc := suite.ratingController // replace CastController with your actual type

	// Register the cast routes
	cc.RegisterRatingRoute(router.Group("/v2"))

	// Perform a GET request to "/v2/cast/get/1"
	w := performRatingRequest(router, "GET", "/v2/rating/get/user/1")

	// Assert the response status code
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// You can add more assertions for other routes if needed
}

func performRatingRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func (suite *RatingTestSuite) TearDownSuite() {
	req, _ := http.NewRequest("DELETE", "/rating/delete/1000/1", nil)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)
}

func (suite *RatingTestSuite) setupRatingRouter() *gin.Engine {
	r := gin.New()

	r.GET("/rating/get/movie/:movie_id", suite.ratingController.GetRatingOfMovie)
	r.GET("/rating/get/user/:user_id", suite.ratingController.GetRatingOfUser)
	r.POST("/rating/create", suite.ratingController.CreateRating)
	r.PATCH("/rating/update", suite.ratingController.UpdateRating)
	r.DELETE("/rating/delete/:userId/:movieId", suite.ratingController.DeleteRating)

	return r
}

func TestRatingTestSuite(t *testing.T) {
	suite.Run(t, new(RatingTestSuite))
}
