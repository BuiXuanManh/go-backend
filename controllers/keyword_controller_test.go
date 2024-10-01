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

type KeywordTestSuite struct {
	suite.Suite
	r                 *gin.Engine
	keywordCollection *mongo.Collection
	keywordService    interfaces.KeywordService
	keywordController KeywordController
	setUpDone         bool
}

func (suite *KeywordTestSuite) SetupTest() {
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
	suite.r = suite.setupKeywordRouter()

	suite.keywordCollection = mongoClient.Database(os.Getenv("DB_NAME")).Collection("keywords")
	suite.keywordService = implementations.NewKeywordService(suite.keywordCollection, ctx)
	suite.keywordController = NewKeywordController(suite.keywordService)

	suite.setUpDone = true
}

func (suite *KeywordTestSuite) TestCreateKeywordSuccessfully() {
	body := strings.NewReader(`{
		"movie_id": 1,
		"keyword_list": [
		  "robbery",
		  "detective",
		  "bank",
		  "obsession"]}`)

	req, _ := http.NewRequest("POST", "/keyword/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	assert.Equal(suite.T(), `{"message":"Successful"}`, w.Body.String())
}

func (suite *KeywordTestSuite) TestCreateExistingKeyword() {
	body := strings.NewReader(`{
		"movie_id": 949,
		"keyword_list": [
		  "robbery",
		  "detective",
		  "bank",
		  "obsession"]}`)

	req, _ := http.NewRequest("POST", "/keyword/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	assert.Equal(suite.T(), `{"message":"write exception: write errors: [E11000 duplicate key error collection: moviesDB.keywords index: movie_id_1 dup key: { movie_id: 949 }]"}`, w.Body.String())
}

func (suite *KeywordTestSuite) TestCreateWrongBindedJSONKeyword() {
	body := strings.NewReader(`{UnbindedJSON}`)

	req, _ := http.NewRequest("POST", "/keyword/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"invalid character 'U' looking for beginning of object key string"}`, w.Body.String())
}

func (suite *KeywordTestSuite) TestCreateWrongBodyStructureInputKeyowrd() {
	body := strings.NewReader(`{"test": true}`)

	req, _ := http.NewRequest("POST", "/keyword/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"wrong input structure"}`, w.Body.String())
}



func (suite *KeywordTestSuite) TestGetAnExistingKeyword() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/keyword/get/1", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// Check the response body
	assert.JSONEq(suite.T(), `{
		"movie_id": 1,
		"keyword_list": [
		  "robbery",
		  "detective",
		  "bank",
		  "obsession"]}`, w.Body.String())
}

func (suite *KeywordTestSuite) TestGetANonExistingKeyword() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/keyword/get/4", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"mongo: no documents in result"}`, w.Body.String())
}

func (suite *KeywordTestSuite) TestGetInvalidMovieIdKeyword() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/keyword/get/fsaf", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"strconv.Atoi: parsing \"fsaf\": invalid syntax"}`, w.Body.String())
}

func (suite *KeywordTestSuite) TestNotFoundKeyword() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/keyword/get/", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	// Check the response body
	assert.Equal(suite.T(), "404 page not found", w.Body.String())
}

func (suite *KeywordTestSuite) TestUpdateUnbindJSONKeyword() {
	body := strings.NewReader(`{UnbindedJSON}`)

	req, _ := http.NewRequest("PATCH", "/keyword/update", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"invalid character 'U' looking for beginning of object key string"}`, w.Body.String())
}

func (suite *KeywordTestSuite) TestUpdateWrongStructureKeyword() {
	body := strings.NewReader(`{"test": true}`)

	req, _ := http.NewRequest("PATCH", "/keyword/update", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"wrong input structure"}`, w.Body.String())
}

func (suite *KeywordTestSuite) TestUpdateNonExistentKeyword() {
	body := strings.NewReader(`{
		"movie_id": 20000,
		"keyword_list": ["comedy"]
	  }`)
	// Create a fake HTTP request
	req, _ := http.NewRequest("PATCH", "/keyword/update", body)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"no matched document found for update"}`, w.Body.String())
}

func (suite *KeywordTestSuite) TestUpdateSuccessfulKeyword() {
	body := strings.NewReader(`{
		"movie_id": 1,
		"keyword_list": ["comedy"]
	  }`)

	req, _ := http.NewRequest("PATCH", "/keyword/update", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	assert.Equal(suite.T(), `{"message":"Successful"}`, w.Body.String())
}

func (suite *KeywordTestSuite) TestDeleteInvalidMovieIdKeyword() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("DELETE", "/keyword/delete/fsaf", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"Invalid movie id"}`, w.Body.String())
}

func (suite *KeywordTestSuite) TestDeleteNonExistentKeyword() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("DELETE", "/keyword/delete/20000", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"no matched document found for delete"}`, w.Body.String())
}

func (suite *KeywordTestSuite) TestRegisterKeywordRouteSuccessfully() {
	// Create a new gin router
	router := gin.New()

	// Create an instance of your CastController
	cc := suite.keywordController // replace CastController with your actual type

	// Register the cast routes
	cc.RegisterKeywordRoute(router.Group("/v2"))

	// Perform a GET request to "/v2/cast/get/1"
	w := performKeywordRequest(router, "GET", "/v2/keyword/get/949")

	// Assert the response status code
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// You can add more assertions for other routes if needed
}

func performKeywordRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func (suite *KeywordTestSuite) TearDownSuite() {
	req, _ := http.NewRequest("DELETE", "/keyword/delete/1", nil)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)
}

func (suite *KeywordTestSuite) setupKeywordRouter() *gin.Engine {
	r := gin.New()
	r.GET("/keyword/get/:id", suite.keywordController.GetKeyword)
	r.POST("/keyword/create", suite.keywordController.CreateKeyword)
	r.PATCH("/keyword/update", suite.keywordController.UpdateKeyword)
	r.DELETE("/keyword/delete/:id", suite.keywordController.DeleteKeyword)
	return r
}

func TestKeywordTestSuite(t *testing.T) {
	suite.Run(t, new(KeywordTestSuite))
}
