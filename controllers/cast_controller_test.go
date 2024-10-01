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

var (
	DEFAULT_PORT = "8080"
	TEST_STRING = "go test cast_controller_test.go cast_controller.go -coverprofile=coverage.out & go tool cover -html=coverage.out"
)

type CastTestSuite struct{
	suite.Suite
	r *gin.Engine
	castCollection *mongo.Collection
	castService interfaces.CastService
	castController CastController
	setUpDone bool
}

func (suite *CastTestSuite) SetupTest() {
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

    suite.castCollection = mongoClient.Database(os.Getenv("DB_NAME")).Collection("casts")
	suite.castService = implementations.NewCastService(suite.castCollection, ctx)
	suite.castController = NewCastController(suite.castService)

	suite.setUpDone = true
}

func (suite *CastTestSuite) TestCreateCastSuccessfully(){
	body := strings.NewReader(`{
		"movie_id": 1,
		"cast": [
		  "{'cast_id': 25, 'character': 'Lt. Vincent Hanna', 'credit_id': '52fe4292c3a36847f80291f5', 'gender': 2, 'id': 1158, 'name': 'Al Pacino', 'order': 0, 'profile_path': '/ks7Ba8x9fJUlP9decBr6Dh5mThX.jpg'}"
		]
	  }`)

	req, _ := http.NewRequest("POST", "/cast/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	assert.Equal(suite.T(), `{"message":"Successful"}`, w.Body.String())
}

func (suite *CastTestSuite) TestCreateWrongBindedJSONCast(){
	body := strings.NewReader(`{UnbindedJSON}`)

	req, _ := http.NewRequest("POST", "/cast/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"invalid character 'U' looking for beginning of object key string"}`, w.Body.String())
}

func (suite *CastTestSuite) TestCreateWrongBodyStructureInputCast(){
	body := strings.NewReader(`{"test": true}`)

	req, _ := http.NewRequest("POST", "/cast/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"wrong input structure"}`, w.Body.String())
}

func (suite *CastTestSuite) TestCreateExistingMovieIdCast(){
	body := strings.NewReader(`{
		"movie_id": 1,
		"cast": [
		  "{'cast_id': 25, 'character': 'Lt. Vincent Hanna', 'credit_id': '52fe4292c3a36847f80291f5', 'gender': 2, 'id': 1158, 'name': 'Al Pacino', 'order': 0, 'profile_path': '/ks7Ba8x9fJUlP9decBr6Dh5mThX.jpg'}"
		]
	  }`)

	req, _ := http.NewRequest("POST", "/cast/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	assert.Equal(suite.T(), `{"message":"write exception: write errors: [E11000 duplicate key error collection: moviesDB.casts index: movie_id_1 dup key: { movie_id: 1 }]"}`, w.Body.String())
}

// All methods that begin with "Test" are run as tests within a suite.
func (suite *CastTestSuite) TestGetAnExistingCast() {	
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/cast/get/451", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *CastTestSuite) TestGetANonExistingCast() {	
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/cast/get/4", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"mongo: no documents in result"}`, w.Body.String())
}

func (suite *CastTestSuite) TestGetCastInvalidMovieId() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/cast/get/fsaf", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"strconv.Atoi: parsing \"fsaf\": invalid syntax"}`, w.Body.String())
}

func (suite *CastTestSuite) TestNotFoundCast() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/cast/get/", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	// Check the response body
	assert.Equal(suite.T(), "404 page not found", w.Body.String())
}

func (suite *CastTestSuite) TestUpdateUnbindJSONCast(){
	body := strings.NewReader(`{UnbindedJSON}`)

	req, _ := http.NewRequest("PATCH", "/cast/update", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"invalid character 'U' looking for beginning of object key string"}`, w.Body.String())
}

func (suite *CastTestSuite) TestUpdateWrongStructureCast(){
	body := strings.NewReader(`{"test": true}`)

	req, _ := http.NewRequest("PATCH", "/cast/update", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"wrong input structure"}`, w.Body.String())
}

func (suite *CastTestSuite) TestUpdateNonExistentCast(){
	body := strings.NewReader(`{
		"movie_id": 20000,
		"cast": [
		  "{'cast_id': 25, 'character': 'Lt. Vincent Hanna', 'credit_id': '52fe4292c3a36847f80291f5', 'gender': 2, 'id': 1158, 'name': 'Al Pacino', 'order': 0, 'profile_path': '/ks7Ba8x9fJUlP9decBr6Dh5mThX.jpg'}"
		]
	  }`)
	// Create a fake HTTP request
	req, _ := http.NewRequest("PATCH", "/cast/update", body)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"no matched document found for update"}`, w.Body.String())
}

func (suite *CastTestSuite) TestUpdateSuccessfulCast(){
	body := strings.NewReader(`{
		"movie_id": 1,
		"cast": [
		  "{'cast_id': 1, 'character': 'Lt. Vincent Hanna', 'credit_id': '52fe4292c3a36847f80291f5', 'gender': 2, 'id': 1, 'name': 'Al Pacino', 'order': 0, 'profile_path': '/ks7Ba8x9fJUlP9decBr6Dh5mThX.jpg'}"
		]
	  }`)

	req, _ := http.NewRequest("PATCH", "/cast/update", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	assert.Equal(suite.T(), `{"message":"Successful"}`, w.Body.String())
}

func (suite *CastTestSuite) TestDeleteInvalidMovieIdCast(){
	// Create a fake HTTP request
	req, _ := http.NewRequest("DELETE", "/cast/delete/fsaf", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"Invalid movie id"}`, w.Body.String())
}

func (suite *CastTestSuite) TestDeleteNonExistentCast(){
	// Create a fake HTTP request
	req, _ := http.NewRequest("DELETE", "/cast/delete/20000", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"no matched document found for delete"}`, w.Body.String())
}

func (suite *CastTestSuite) TestRegisterRouteSuccessfully(){
	// Create a new gin router
	router := gin.New()

	// Create an instance of your CastController
	cc := suite.castController // replace CastController with your actual type

	// Register the cast routes
	cc.RegisterCastRoute(router.Group("/v2"))

	// Perform a GET request to "/v2/cast/get/1"
	w := performRequest(router, "GET", "/v2/cast/get/949")

	// Assert the response status code
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// You can add more assertions for other routes if needed
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}


// func (suite *CastTestSuite) TestDeleteSuccessfulCast(){
// 	req, _ := http.NewRequest("DELETE", "/cast/delete/1", nil)

// 	w := httptest.NewRecorder()

// 	suite.r.ServeHTTP(w, req)

// 	assert.Equal(suite.T(), http.StatusOK, w.Code)

// 	assert.Equal(suite.T(), `{"message":"Successful"}`, w.Body.String())
// }


func (suite *CastTestSuite) setupRouter() *gin.Engine {
	r := gin.New()
	r.GET("/cast/get/:id", suite.castController.GetCast)
	r.POST("/cast/create", suite.castController.CreateCast)
	r.PATCH("/cast/update", suite.castController.UpdateCast)
	r.DELETE("/cast/delete/:id", suite.castController.DeleteCast)
	return r
}


// TearDownSuite is called once after all tests in the suite have been run
func (suite *CastTestSuite) TearDownSuite() {
	req, _ := http.NewRequest("DELETE", "/cast/delete/1", nil)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)
}

func TestCastTestSuite (t *testing.T){
	suite.Run(t, new(CastTestSuite))
}