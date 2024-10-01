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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MovieDiscussionTestSuite struct {
	suite.Suite
	r                         *gin.Engine
	movieDiscussionCollection *mongo.Collection
	movieDiscussionService    interfaces.MovieDiscussionServices
	movieDiscussionController MovieDiscussionController
	discussionID              string
	setUpDone                 bool
}

func (suite *MovieDiscussionTestSuite) SetupTest() {
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
	suite.r = suite.setupMovieDiscussionRouter()

	suite.movieDiscussionCollection = mongoClient.Database(os.Getenv("DB_NAME")).Collection("discussion_movie")

	suite.movieDiscussionService = implementations.NewMovieDiscussionServices(suite.movieDiscussionCollection, ctx)

	suite.movieDiscussionController = NewMovieDiscussionController(suite.movieDiscussionService)

	suite.setUpDone = true
}

func (suite *MovieDiscussionTestSuite) TestCreateMovieDiscussionSuccessfully() {
	id := primitive.NewObjectID().Hex()

	body := strings.NewReader(`{
		"_id": "` + id + `",
		"movie_id": 949,
		"subject": "Thoughts on the Ending",
		"status": true,
		"discussion_part": [
			{
				"user_id": 672,
				"part_id": 0,
				"name": "4nh3k",
				"profile_path": "https://ui-avatars.com/api/?background=01d277&name=4nh3k&color=fff&size=256",
				"title": "",
				"timestamp":  "2023-12-25T14:56:26.142Z",
				"description": "Test",
				"is_reply_of": null
			  }
		]
	}`)

	req, _ := http.NewRequest("POST", "/movieDiscussion/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	assert.Equal(suite.T(), `{"message":"Successful"}`, w.Body.String())

	suite.discussionID = id

}

func (suite *MovieDiscussionTestSuite) TestCreateWrongBindedJSONMovieDiscussion() {
	body := strings.NewReader(`{UnbindedJSON}`)

	req, _ := http.NewRequest("POST", "/movieDiscussion/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"invalid character 'U' looking for beginning of object key string"}`, w.Body.String())
}

func (suite *MovieDiscussionTestSuite) TestGetAnExistingMovieDiscussion() {
	id := suite.discussionID

	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/movieDiscussion/get/"+id, nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// Check the response body
	assert.JSONEq(suite.T(), `{
		"_id": "` + id + `",
		"discussion_part": [
			{
				"description": "Test",
				"user_id": 672,
				"name": "4nh3k",
				"part_id": 0,
				"profile_path": "https://ui-avatars.com/api/?background=01d277&name=4nh3k&color=fff&size=256",
				"timestamp":  "2023-12-25T14:56:26.142Z",
				"title": ""
			  }
		],
		"movie_id": 949,
		"status": true,
		"subject": "Thoughts on the Ending"
	}`,
		 w.Body.String())
}

func (suite *MovieDiscussionTestSuite) TestGetANonExistingMovieDiscussion() {
	// Create a fake HTTP request
	id := primitive.NewObjectID().Hex()

	req, _ := http.NewRequest("GET", "/movieDiscussion/get/"+id, nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"mongo: no documents in result"}`, w.Body.String())
}

func (suite *MovieDiscussionTestSuite) TestGetInvalidMovieIdMovieDiscussion() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/movieDiscussion/get/fsaf", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"the provided hex string is not a valid ObjectID"}`, w.Body.String())
}

func (suite *MovieDiscussionTestSuite) TestGetANonExistingMovieByMovieDiscussion() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/movieDiscussion/getByMovie/2000001", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"mongo: no documents in result"}`, w.Body.String())
}
func (suite *MovieDiscussionTestSuite) TestGetInvalidMovieIdByMovieDiscussion() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/movieDiscussion/getByMovie/fsaf", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"strconv.Atoi: parsing \"fsaf\": invalid syntax"}`, w.Body.String())
}

func (suite *MovieDiscussionTestSuite) TestGetANonExistingMovieByUserDiscussion() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/movieDiscussion/getByUser/2000001", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"mongo: no documents in result"}`, w.Body.String())
}
func (suite *MovieDiscussionTestSuite) TestGetInvalidMovieIdByUserDiscussion() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/movieDiscussion/getByUser/fsaf", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"strconv.Atoi: parsing \"fsaf\": invalid syntax"}`, w.Body.String())
}

func (suite *MovieDiscussionTestSuite) TestUpdateUnbindJSONMovieDiscussion() {
	body := strings.NewReader(`{UnbindedJSON}`)

	req, _ := http.NewRequest("PATCH", "/movieDiscussion/update", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"invalid character 'U' looking for beginning of object key string"}`, w.Body.String())
}

func (suite *MovieDiscussionTestSuite) TestUpdateNonExistentMovieDiscussion() {
	id := primitive.NewObjectID().Hex()

	body := strings.NewReader(`{
		"_id": "` + id + `",
		"movie_id": 123,
		"subject": "Thoughts on the Ending",
		"status": true,
		"discussion_part": []
	  }`)
	// Create a fake HTTP request
	req, _ := http.NewRequest("PATCH", "/movieDiscussion/update", body)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"no matched document found for update"}`, w.Body.String())
}

func (suite *MovieDiscussionTestSuite) TestUpdateSuccessfulMovieDiscussion() {
	id := suite.discussionID
	body := strings.NewReader(`{
		"_id": "` + id + `",
		"movie_id": 123,
		"subject": "Thoughts on the Ending",
		"status": true,
		"discussion_part": []
	  }`)

	req, _ := http.NewRequest("PATCH", "/movieDiscussion/update", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	assert.Equal(suite.T(), `{"message":"Successful"}`, w.Body.String())
}

func (suite *MovieDiscussionTestSuite) TestDeleteInvalidMovieIdMovieDiscussion() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("DELETE", "/movieDiscussion/delete/fsaf", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"Invalid movie id"}`, w.Body.String())
}

func (suite *MovieDiscussionTestSuite) TestDeleteNonExistentMovieDiscussion() {
	// Create a fake HTTP request
	id := primitive.NewObjectID().Hex()
	req, _ := http.NewRequest("DELETE", "/movieDiscussion/delete/"+id, nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"no matched document found for delete"}`, w.Body.String())
}

func (suite *MovieDiscussionTestSuite) TestCreatePartInInvalidDiscussion() {
	body := strings.NewReader(`{}`)

	req, _ := http.NewRequest("PATCH", "/movieDiscussion/create/part/fsaf", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"the provided hex string is not a valid ObjectID"}`, w.Body.String())
}

func (suite *MovieDiscussionTestSuite) TestCreateWrongBindedJSONMovieDiscussionPart() {
	body := strings.NewReader(`{UnbindedJSON}`)
	id := suite.discussionID
	req, _ := http.NewRequest("PATCH", "/movieDiscussion/create/part/"+id, body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"invalid character 'U' looking for beginning of object key string"}`, w.Body.String())
}

func (suite *MovieDiscussionTestSuite) TestUpdateMovieDiscussionPartInvalidDiscussionMovieId() {
	body := strings.NewReader(`{}`)

	req, _ := http.NewRequest("PATCH", "/movieDiscussion/update/part/fsaf/1", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"Invalid discussion or part id"}`, w.Body.String())
}

func (suite *MovieDiscussionTestSuite) TestUpdateWrongBindedJSONMovieDiscussionPart() {
	body := strings.NewReader(`{UnbindedJSON}`)

	id := suite.discussionID

	req, _ := http.NewRequest("PATCH", "/movieDiscussion/update/part/"+id+"/0", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"invalid character 'U' looking for beginning of object key string"}`, w.Body.String())
}

func (suite *MovieDiscussionTestSuite) TestDeleteInvalidDiscussionPart() {
	body := strings.NewReader(``)

	id := suite.discussionID

	req, _ := http.NewRequest("PATCH", "/movieDiscussion/delete/part/"+id+"/fsaf", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"Invalid discussion or part id"}`, w.Body.String())
}

func (suite *MovieDiscussionTestSuite) TestGetMovieDiscussionInPageFromInvalidDiscussionId() {
	body := strings.NewReader(`{}`)

	req, _ := http.NewRequest("GET", "/movieDiscussion/get/discussion/page/fsaf", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"Invalid page number"}`, w.Body.String())
}

func (suite *MovieDiscussionTestSuite) TestGetMovieDiscussionInPageFromNonExistentPage() {
	body := strings.NewReader(`{}`)

	req, _ := http.NewRequest("GET", "/movieDiscussion/get/discussion/page/2000000", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	assert.Equal(suite.T(), `{"message":"documents not found"}`, w.Body.String())
}

func (suite *MovieDiscussionTestSuite) TestRegisterMovieDiscussionRouteSuccessfully() {
	// Create a new gin router
	router := gin.New()

	// Create an instance of your CastController
	cc := suite.movieDiscussionController // replace CastController with your actual type

	// Register the cast routes
	cc.RegisterMovieDiscussionRoute(router.Group("/v2"))
	id := suite.discussionID
	// Perform a GET request to "/v2/cast/get/1"
	w := performMovieDiscussionRequest(router, "GET", "/v2/movieDiscussion/get/"+id)

	// Assert the response status code
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// You can add more assertions for other routes if needed
}

func performMovieDiscussionRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func (suite *MovieDiscussionTestSuite) TearDownSuite() {

	id := suite.discussionID

	req, _ := http.NewRequest("DELETE", "/movieDiscussion/delete/"+id, nil)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)
}

func (suite *MovieDiscussionTestSuite) setupMovieDiscussionRouter() *gin.Engine {
	r := gin.New()
	r.GET("/movieDiscussion/get/:id", suite.movieDiscussionController.GetMovieDiscussion)

	r.POST("/movieDiscussion/create", suite.movieDiscussionController.CreateMovieDiscussion)

	r.PATCH("/movieDiscussion/update", suite.movieDiscussionController.UpdateMovieDiscussion)

	r.DELETE("/movieDiscussion/delete/:id", suite.movieDiscussionController.DeleteMovieDiscussion)

	r.GET("/movieDiscussion/getByMovie/:movie_id", suite.movieDiscussionController.GetMovieDiscussionsByMovieId)

	r.GET("/movieDiscussion/getByUser/:user_id", suite.movieDiscussionController.GetMovieDiscussionsByUserId)

	r.PATCH("/movieDiscussion/create/part/:discussion_id", suite.movieDiscussionController.CreateMovieDiscussionPart)

	r.GET("/movieDiscussion/get/discussion/page/:pageNumber", suite.movieDiscussionController.GetMovieDiscussionInPage)

	r.PATCH("/movieDiscussion/update/part/:discussion_id/:part_id", suite.movieDiscussionController.UpdateMovieDiscussionPart)

	r.PATCH("/movieDiscussion/delete/part/:discussion_id/:part_id", suite.movieDiscussionController.DeleteMovieDiscussionPart)
	return r
}

func TestMovieDiscussionTestSuite(t *testing.T) {
	suite.Run(t, new(MovieDiscussionTestSuite))
}
