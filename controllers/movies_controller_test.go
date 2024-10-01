package controllers

import (
	"context"
	"log"
	"gotest/constants"
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

var (
	movieJSONForTesting          = constants.MOVIE_JSON_FOR_TESTING
	movieJSONForTestingPopular   = constants.MOVIE_JSON_FOR_TESTING_POPULAR
	movieDuplicateJSONForTesting = constants.MOVIE_DUPLICATE_FOR_TESTING
	movieNonExistent             = constants.MOVIE_NON_EXISTENT
	movieForEdit                 = constants.MOVIE_EDIT
)

type MovieTestSuite struct {
	suite.Suite
	r               *gin.Engine
	movieCollection *mongo.Collection
	movieService    interfaces.MovieService
	movieController MovieController
	setUpDone       bool
}

func (suite *MovieTestSuite) SetupTest() {
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
	suite.r = suite.setupMovieRouter()

	suite.movieCollection = mongoClient.Database(os.Getenv("DB_NAME")).Collection("movies")
	suite.movieService = implementations.NewMovieService(suite.movieCollection, ctx)
	suite.movieController = NewMovieController(suite.movieService)

	suite.setUpDone = true
}

func (suite *MovieTestSuite) TestCreateMovieSuccessfully() {
	body := strings.NewReader(movieJSONForTesting)

	req, _ := http.NewRequest("POST", "/movie/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	assert.Equal(suite.T(), `{"message":"Successful"}`, w.Body.String())
}

func (suite *MovieTestSuite) TestCreateExistingMovie() {
	body := strings.NewReader(movieDuplicateJSONForTesting)

	req, _ := http.NewRequest("POST", "/movie/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	assert.Equal(suite.T(), `{"message":"write exception: write errors: [E11000 duplicate key error collection: moviesDB.movies index: id_1 dup key: { id: 949 }]"}`, w.Body.String())
}

func (suite *MovieTestSuite) TestCreateWrongBindedJSONMovie() {
	body := strings.NewReader(`{UnbindedJSON}`)

	req, _ := http.NewRequest("POST", "/movie/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"invalid character 'U' looking for beginning of object key string"}`, w.Body.String())
}

func (suite *MovieTestSuite) TestCreateWrongBodyStructureInputMovie() {
	body := strings.NewReader(`{"test": true}`)

	req, _ := http.NewRequest("POST", "/movie/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"wrong input structure"}`, w.Body.String())
}

func (suite *MovieTestSuite) TestCreateWrongBindedJSONMovies() {
	body := strings.NewReader(`{UnbindedJSON}`)

	req, _ := http.NewRequest("POST", "/movie/createMany", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"invalid character 'U' looking for beginning of object key string"}`, w.Body.String())
}

func (suite *MovieTestSuite) TestGetAnExistingMovie() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/movie/get/949", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// Check the response body
	assert.JSONEq(suite.T(), movieDuplicateJSONForTesting, w.Body.String())
}

func (suite *MovieTestSuite) TestGetANonExistingMovie() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/movie/get/4", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"mongo: no documents in result"}`, w.Body.String())
}

func (suite *MovieTestSuite) TestGetInvalidMovieIdMovie() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/movie/get/fsaf", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"strconv.Atoi: parsing \"fsaf\": invalid syntax"}`, w.Body.String())
}

func (suite *MovieTestSuite) TestNotFoundMovie() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/movie/get/", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	// Check the response body
	assert.Equal(suite.T(), "404 page not found", w.Body.String())
}

func (suite *MovieTestSuite) TestGetPopularMoviesSuccessfully() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/movie/get/popular/1", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// Check the response body
	// assert.Equal(suite.T(), movieJSONForTestingPopular, w.Body.String())
}

func (suite *MovieTestSuite) TestGetInvalidNumberOfMoviesPopularMovies() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/movie/get/popular/-1", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"Invalid number of movies"}`, w.Body.String())
}

func (suite *MovieTestSuite) TestGetInvalidPageNumberMovies() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/movie/get/page/-1", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Check the response body
	assert.JSONEq(suite.T(), `{"message":"Invalid page number"}`, w.Body.String())
}

func (suite *MovieTestSuite) TestGetBadGatewayPageMovies() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/movie/get/page/20000", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.JSONEq(suite.T(), `{"message":"documents not found"}`, w.Body.String())
}

func (suite *MovieTestSuite) TestSearchInvalidPageNumberMovies() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/movie/search/cutthro/-1", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Check the response body
	assert.JSONEq(suite.T(), `{"message":"Invalid page number"}`, w.Body.String())
}

func (suite *MovieTestSuite) TestSearchBadGatewayPageMovies() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/movie/search/cutthroat/20000", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.JSONEq(suite.T(), `{"message":"documents not found"}`, w.Body.String())
}

func (suite *MovieTestSuite) TestUpdateUnbindJSONMovie() {
	body := strings.NewReader(`{UnbindedJSON}`)

	req, _ := http.NewRequest("PATCH", "/movie/update", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"invalid character 'U' looking for beginning of object key string"}`, w.Body.String())
}

func (suite *MovieTestSuite) TestUpdateWrongStructureMovie() {
	body := strings.NewReader(`{"test": true}`)

	req, _ := http.NewRequest("PATCH", "/movie/update", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"wrong input structure"}`, w.Body.String())
}

func (suite *MovieTestSuite) TestUpdateNonExistentMovie() {
	body := strings.NewReader(movieNonExistent)
	// Create a fake HTTP request
	req, _ := http.NewRequest("PATCH", "/movie/update", body)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"no matched document found for update"}`, w.Body.String())
}

func (suite *MovieTestSuite) TestUpdateSuccessfulMovie() {
	body := strings.NewReader(movieForEdit)

	req, _ := http.NewRequest("PATCH", "/movie/update", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	assert.Equal(suite.T(), `{"message":"Successful"}`, w.Body.String())
}

func (suite *MovieTestSuite) TestDeleteInvalidMovieIdMovie() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("DELETE", "/movie/delete/fsaf", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"Invalid movie id"}`, w.Body.String())
}

func (suite *MovieTestSuite) TestDeleteNonExistentMovie() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("DELETE", "/movie/delete/20000", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"no matched document found for delete"}`, w.Body.String())
}

func (suite *MovieTestSuite) TestRegisterMovieRouteSuccessfully() {
	// Create a new gin router
	router := gin.New()

	// Create an instance of your CastController
	cc := suite.movieController // replace CastController with your actual type

	// Register the cast routes
	cc.RegisterMovieRoute(router.Group("/v2"))

	// Perform a GET request to "/v2/movie/get/1"
	w := performMovieRequest(router, "GET", "/v2/movie/get/949")

	// Assert the response status code
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// You can add more assertions for other routes if needed
}

func performMovieRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func (suite *MovieTestSuite) TearDownSuite() {
	req, _ := http.NewRequest("DELETE", "/movie/delete/1", nil)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)
}

func (suite *MovieTestSuite) setupMovieRouter() *gin.Engine {
	r := gin.New()
	r.POST("/movie/create", suite.movieController.CreateMovie)
	r.POST("/movie/createMany", suite.movieController.CreateMovies)

	r.GET("/movie/get/:id", suite.movieController.GetMovie)
	r.GET("/movie/get/popular/:numberOfMovies", suite.movieController.GetPopularMovies)
	r.GET("/movie/get/page/:pageNumber", suite.movieController.GetMoviesInPage)
	r.GET("/movie/search/:search_word/:pageNumber", suite.movieController.SearchMovie)

	r.PATCH("/movie/update", suite.movieController.UpdateMovie)
	r.DELETE("/movie/delete/:id", suite.movieController.DeleteMovie)
	return r
}

func TestMovieTestSuite(t *testing.T) {
	suite.Run(t, new(MovieTestSuite))
}
