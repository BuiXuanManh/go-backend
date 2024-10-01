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

type UserReviewTestSuite struct {
	suite.Suite
	r                    *gin.Engine
	UserReviewCollection *mongo.Collection
	UserReviewService    interfaces.UserReviewService
	UserReviewController UserReviewController
	setUpDone            bool
}

func (suite *UserReviewTestSuite) SetupTest() {
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

	suite.UserReviewCollection = mongoClient.Database(os.Getenv("DB_NAME")).Collection("user_reviews")
	suite.UserReviewService = implementations.NewUserReviewService(suite.UserReviewCollection, ctx)
	suite.UserReviewController = NewUserReviewController(suite.UserReviewService)

	suite.setUpDone = true
}

func (suite *UserReviewTestSuite) TestCreateUserReviewSuccessfully() {
	body := strings.NewReader(`{
		"user_id": 672,
		"movie_id": 2048,
		"username": "4nh3k",
		"picture_profile": "https://picsum.photos/672",
		"comment": "Incredible film! The storyline kept me on the edge of my seat, and the performances were outstanding."
		}`)

	req, _ := http.NewRequest("POST", "/userReview/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	assert.Equal(suite.T(), `{"message":"Successful"}`, w.Body.String())
}

func (suite *UserReviewTestSuite) TestCreateWrongBindedJSONUserReview() {
	body := strings.NewReader(`{UnbindedJSON}`)

	req, _ := http.NewRequest("POST", "/userReview/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"invalid character 'U' looking for beginning of object key string"}`, w.Body.String())
}

func (suite *UserReviewTestSuite) TestCreateExistingUserIdMovieIdReview() {
	body := strings.NewReader(`{
		"user_id": 672,
		"movie_id": 1024,
		"username": "4nh3k",
		"picture_profile": "https://picsum.photos/672",
		"comment": "Incredible film! The storyline kept me on the edge of my seat, and the performances were outstanding."
	}`)

	req, _ := http.NewRequest("POST", "/userReview/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	assert.Equal(suite.T(), `{"message":"write exception: write errors: [E11000 duplicate key error collection: moviesDB.user_reviews index: movie_id_1_user_id_1 dup key: { movie_id: 1024, user_id: 672 }]"}`, w.Body.String())
}

// ------------------------------------------------------------
func (suite *UserReviewTestSuite) TestGetANonExistingUserReview() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/userReview/get/user/20000", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"mongo: no documents in result"}`, w.Body.String())
}

func (suite *UserReviewTestSuite) TestGetUserReviewInvalidMovieId() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/userReview/get/user/fsaf", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"strconv.Atoi: parsing \"fsaf\": invalid syntax"}`, w.Body.String())
}

func (suite *UserReviewTestSuite) TestNotFoundUserReview() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/userReview/get/user", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	// Check the response body
	assert.Equal(suite.T(), "404 page not found", w.Body.String())
}

func (suite *UserReviewTestSuite) TestGetANonExistingMovieReview() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/userReview/get/movie/20000", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"mongo: no documents in result"}`, w.Body.String())
}

func (suite *UserReviewTestSuite) TestGetMovieReviewInvalidMovieId() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/userReview/get/movie/fsaf", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"strconv.Atoi: parsing \"fsaf\": invalid syntax"}`, w.Body.String())
}

func (suite *UserReviewTestSuite) TestNotFoundMovieReview() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/userReview/get/movie", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	// Check the response body
	assert.Equal(suite.T(), "404 page not found", w.Body.String())
}

// ---------------------------------------------------------

func (suite *UserReviewTestSuite) TestUpdateUnbindJSONUserReview() {
	body := strings.NewReader(`{UnbindedJSON}`)

	req, _ := http.NewRequest("PATCH", "/userReview/update", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"invalid character 'U' looking for beginning of object key string"}`, w.Body.String())
}

func (suite *UserReviewTestSuite) TestUpdateNonExistentUserReview() {
	body := strings.NewReader(`{
		"user_id": 1,
		"movie_id": 1,
		"username": "4nh3k",
		"picture_profile": "https://picsum.photos/672",
		"comment": "Incredible film! The storyline kept me on the edge of my seat, and the performances were outstanding."
		}`)

	// Create a fake HTTP request
	req, _ := http.NewRequest("PATCH", "/userReview/update", body)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"no matched document found for update"}`, w.Body.String())
}

func (suite *UserReviewTestSuite) TestUpdateSuccessfulUserReview() {
	body := strings.NewReader(`{
		"user_id": 672,
		"movie_id": 1024,
		"username": "4nh3k",
		"picture_profile": "https://picsum.photos/672",
		"comment": "Great movie! I loved the plot twists and the acting was superb."
		}`)

	req, _ := http.NewRequest("PATCH", "/userReview/update", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	assert.Equal(suite.T(), `{"message":"Successful"}`, w.Body.String())
}

func (suite *UserReviewTestSuite) TestDeleteInvalidMovieIdUserReview() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("DELETE", "/userReview/delete/fsaf/fsaf", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"Invalid movie id"}`, w.Body.String())
}

func (suite *UserReviewTestSuite) TestDeleteNonExistentUserReview() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("DELETE", "/userReview/delete/20000/20000", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"no matched document found for delete"}`, w.Body.String())
}

func (suite *UserReviewTestSuite) TestRegisterRouteSuccessfully() {
	// Create a new gin router
	router := gin.New()

	// Create an instance of your UserReviewController
	cc := suite.UserReviewController // replace UserReviewController with your actual type

	// Register the UserReview routes
	cc.RegisterUserReviewRoute(router.Group("/v2"))

	// Perform a GET request to "/v2/UserReview/get/1"
	w := performUserReviewRequest(router, "GET", "/v2/userReview/get/user/672")

	// Assert the response status code
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// You can add more assertions for other routes if needed
}

func performUserReviewRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// func (suite *UserReviewTestSuite) TestDeleteSuccessfulUserReview(){
// 	req, _ := http.NewRequest("DELETE", "/UserReview/delete/1", nil)

// 	w := httptest.NewRecorder()

// 	suite.r.ServeHTTP(w, req)

// 	assert.Equal(suite.T(), http.StatusOK, w.Code)

// 	assert.Equal(suite.T(), `{"message":"Successful"}`, w.Body.String())
// }

func (suite *UserReviewTestSuite) setupRouter() *gin.Engine {
	r := gin.New()
	r.GET("userReview/get/user/:user_id", suite.UserReviewController.GetUserReviewByUserId)

	r.GET("userReview/get/movie/:movie_id", suite.UserReviewController.GetUserReviewByMovieId)

	r.POST("userReview/create", suite.UserReviewController.CreateUserReview)

	r.PATCH("userReview/update", suite.UserReviewController.UpdateUserReview)

	r.DELETE("userReview/delete/:user_id/:movie_id", suite.UserReviewController.DeleteUserReview)
	return r
}

// TearDownSuite is called once after all tests in the suite have been run
func (suite *UserReviewTestSuite) TearDownSuite() {
	req, _ := http.NewRequest("DELETE", "/userReview/delete/672/2048", nil)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)
}

func TestUserReviewTestSuite(t *testing.T) {
	suite.Run(t, new(UserReviewTestSuite))
}
