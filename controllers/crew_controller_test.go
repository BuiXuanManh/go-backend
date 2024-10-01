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

type CrewTestSuite struct {
	suite.Suite
	r              *gin.Engine
	crewCollection *mongo.Collection
	crewService    interfaces.CrewService
	crewController CrewController
	setUpDone      bool
}

func (suite *CrewTestSuite) SetupTest() {
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
	suite.r = suite.setupCrewRouter()

	suite.crewCollection = mongoClient.Database(os.Getenv("DB_NAME")).Collection("crews")
	suite.crewService = implementations.NewCrewService(suite.crewCollection, ctx)
	suite.crewController = NewCrewController(suite.crewService)

	suite.setUpDone = true
}

func (suite *CrewTestSuite) TestCreateCrewSuccessfully() {
	body := strings.NewReader(`{
		"movie_id": 1,
		"crew": [
		  "{'credit_id': '52fe4292c3a36847f802916d', 'department': 'Directing', 'gender': 2, 'id': 638, 'job': 'Director', 'name': 'Michael Mann', 'profile_path': '/nKmUpRpuQIsYubR7vIxVKhkbaTW.jpg'}"
		]
	  }`)

	req, _ := http.NewRequest("POST", "/crew/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	assert.Equal(suite.T(), `{"message":"Successful"}`, w.Body.String())
}

func (suite *CrewTestSuite) TestCreateWrongBindedJSONCrew() {
	body := strings.NewReader(`{UnbindedJSON}`)

	req, _ := http.NewRequest("POST", "/crew/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"invalid character 'U' looking for beginning of object key string"}`, w.Body.String())
}

func (suite *CrewTestSuite) TestCreateWrongBodyStructureInputCrew() {
	body := strings.NewReader(`{"test": true}`)

	req, _ := http.NewRequest("POST", "/crew/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"wrong input structure"}`, w.Body.String())
}

func (suite *CrewTestSuite) TestCreateExistingMovieIdCrew() {
	body := strings.NewReader(`{
		"movie_id": 1,
		"crew": [
		  "{'credit_id': '52fe4292c3a36847f802916d', 'department': 'Directing', 'gender': 2, 'id': 638, 'job': 'Director', 'name': 'Michael Mann', 'profile_path': '/nKmUpRpuQIsYubR7vIxVKhkbaTW.jpg'}"
		]
	  }`)

	req, _ := http.NewRequest("POST", "/crew/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	assert.Equal(suite.T(), `{"message":"write exception: write errors: [E11000 duplicate key error collection: moviesDB.crews index: movie_id_1 dup key: { movie_id: 1 }]"}`, w.Body.String())
}

func (suite *CrewTestSuite) TestGetAnExistingCrew() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/crew/get/1", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// Check the response body
	assert.JSONEq(suite.T(), `{
		"movie_id": 1,
		"crew": [
		  "{'credit_id': '52fe4292c3a36847f802916d', 'department': 'Directing', 'gender': 2, 'id': 638, 'job': 'Director', 'name': 'Michael Mann', 'profile_path': '/nKmUpRpuQIsYubR7vIxVKhkbaTW.jpg'}"
		]
	  }`, w.Body.String())
}

func (suite *CrewTestSuite) TestGetANonExistingCrew() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/crew/get/4", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"mongo: no documents in result"}`, w.Body.String())
}

func (suite *CrewTestSuite) TestGetInvalidMovieIdCrew() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/crew/get/fsaf", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"strconv.Atoi: parsing \"fsaf\": invalid syntax"}`, w.Body.String())
}

func (suite *CrewTestSuite) TestNotFoundCrew() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/crew/get/", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	// Check the response body
	assert.Equal(suite.T(), "404 page not found", w.Body.String())
}

func (suite *CrewTestSuite) TestUpdateUnbindJSONCrew() {
	body := strings.NewReader(`{UnbindedJSON}`)

	req, _ := http.NewRequest("PATCH", "/crew/update", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"invalid character 'U' looking for beginning of object key string"}`, w.Body.String())
}

func (suite *CrewTestSuite) TestUpdateWrongStructureCrew() {
	body := strings.NewReader(`{"test": true}`)

	req, _ := http.NewRequest("PATCH", "/crew/update", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"wrong input structure"}`, w.Body.String())
}

func (suite *CrewTestSuite) TestUpdateNonExistentCrew() {
	body := strings.NewReader(`{
		"movie_id": 20000,
		"crew": [
		  "{'credit_id': '52fe4292c3a36847f802916d', 'department': 'Directing', 'gender': 2, 'id': 638, 'job': 'Director', 'name': 'Michael Mann', 'profile_path': '/nKmUpRpuQIsYubR7vIxVKhkbaTW.jpg'}"
		]
	  }`)
	// Create a fake HTTP request
	req, _ := http.NewRequest("PATCH", "/crew/update", body)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"no matched document found for update"}`, w.Body.String())
}

func (suite *CrewTestSuite) TestUpdateSuccessfulCrew() {
	body := strings.NewReader(`{
		"movie_id": 1,
		"crew": [
		  "{'credit_id': '52fe4292c3a36847f802916d', 'department': 'Directing', 'gender': 2, 'id': 1, 'job': 'Director', 'name': 'Michael Mann', 'profile_path': '/nKmUpRpuQIsYubR7vIxVKhkbaTW.jpg'}"
		]
	  }`)

	req, _ := http.NewRequest("PATCH", "/crew/update", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	assert.Equal(suite.T(), `{"message":"Successful"}`, w.Body.String())
}

func (suite *CrewTestSuite) TestDeleteInvalidMovieIdCrew() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("DELETE", "/crew/delete/fsaf", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"Invalid movie id"}`, w.Body.String())
}

func (suite *CrewTestSuite) TestDeleteNonExistentCrew() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("DELETE", "/crew/delete/20000", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"no matched document found for delete"}`, w.Body.String())
}

func (suite *CrewTestSuite) TestRegisterCrewRouteSuccessfully() {
	// Create a new gin router
	router := gin.New()

	// Create an instance of your CastController
	cc := suite.crewController // replace CastController with your actual type

	// Register the cast routes
	cc.RegisterCrewRoute(router.Group("/v2"))

	// Perform a GET request to "/v2/cast/get/1"
	w := performCrewRequest(router, "GET", "/v2/crew/get/949")

	// Assert the response status code
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// You can add more assertions for other routes if needed
}

func performCrewRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func (suite *CrewTestSuite) TearDownSuite() {
	req, _ := http.NewRequest("DELETE", "/crew/delete/1", nil)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)
}

func (suite *CrewTestSuite) setupCrewRouter() *gin.Engine {
	r := gin.New()
	r.GET("/crew/get/:id", suite.crewController.GetCrew)
	r.POST("/crew/create", suite.crewController.CreateCrew)
	r.PATCH("/crew/update", suite.crewController.UpdateCrew)
	r.DELETE("/crew/delete/:id", suite.crewController.DeleteCrew)
	return r
}

func TestCrewTestSuite(t *testing.T) {
	suite.Run(t, new(CrewTestSuite))
}
