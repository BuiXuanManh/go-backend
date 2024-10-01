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

type UserTestSuite struct {
	suite.Suite
	r              *gin.Engine
	userCollection *mongo.Collection
	userService    interfaces.UserService
	userController UserController
	setUpDone      bool
}

func (suite *UserTestSuite) SetupTest() {
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
	suite.r = suite.setupUserRouter()

	suite.userCollection = mongoClient.Database(os.Getenv("DB_NAME")).Collection("users")
	suite.userService = implementations.NewUserService(suite.userCollection, ctx)
	suite.userController = NewUserController(suite.userService)

	suite.setUpDone = true
}

func (suite *UserTestSuite) TestCreateUserSuccessfully() {
	body := strings.NewReader(`{
		"id": 20000,
		"username": "test20000",
		"password_hash": "abcxyz",
		"favorite_list": [2,3,5,6,11],
		"recommendation_list": [2,3,5,6,11],
		"watch_list": [2,3,5,6,11],
		"picture_profile": "" }`)

	req, _ := http.NewRequest("POST", "/user/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	assert.Equal(suite.T(), `{"message":"Successful"}`, w.Body.String())
}

func (suite *UserTestSuite) TestCreateExistingUserUsername() {
	body := strings.NewReader(`{
		"id": 10001,
		"username": "test",
		"password_hash": "abcxyz",
		"favorite_list": [2,3,5,6,11],
		"recommendation_list": [2,3,5,6,11],
		"watch_list": [2,3,5,6,11],
		"picture_profile": "" }`)

	req, _ := http.NewRequest("POST", "/user/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	assert.Equal(suite.T(), `{"message":"write exception: write errors: [E11000 duplicate key error collection: moviesDB.users index: username_1 dup key: { username: \"test\" }]"}`, w.Body.String())
}

func (suite *UserTestSuite) TestCreateWrongBindedJSONUser() {
	body := strings.NewReader(`{UnbindedJSON}`)

	req, _ := http.NewRequest("POST", "/user/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"invalid character 'U' looking for beginning of object key string"}`, w.Body.String())
}

func (suite *UserTestSuite) TestCreateWrongBodyStructureInputUser() {
	body := strings.NewReader(`{"test": true}`)

	req, _ := http.NewRequest("POST", "/user/create", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"wrong input structure"}`, w.Body.String())
}

func (suite *UserTestSuite) TestGetAnExistingUser() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/user/get/10001", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// Check the response body
	assert.JSONEq(suite.T(), `{
		"id": 10001,
		"username": "test",
		"password_hash": "abcxyz",
		"favorite_list": [2,3,5,6,11],
		"recommendation_list": [2,3,5,6,11],
		"watch_list": [2,3,5,6,11],
		"picture_profile": "" }`, w.Body.String())
}

func (suite *UserTestSuite) TestGetANonExistingUser() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/user/get/30000", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"mongo: no documents in result"}`, w.Body.String())
}

func (suite *UserTestSuite) TestGetInvalidUserId() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/user/get/fsaf", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"strconv.Atoi: parsing \"fsaf\": invalid syntax"}`, w.Body.String())
}

func (suite *UserTestSuite) TestNotFoundUser() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "/user/get/", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	// Check the response body
	assert.Equal(suite.T(), "404 page not found", w.Body.String())
}

func (suite *UserTestSuite) TestUpdateUnbindJSONUser() {
	body := strings.NewReader(`{UnbindedJSON}`)

	req, _ := http.NewRequest("PATCH", "/user/update", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"invalid character 'U' looking for beginning of object key string"}`, w.Body.String())
}

func (suite *UserTestSuite) TestUpdateWrongStructureKeyword() {
	body := strings.NewReader(`{"test": true}`)

	req, _ := http.NewRequest("PATCH", "/user/update", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"wrong input structure"}`, w.Body.String())
}

func (suite *UserTestSuite) TestUpdateNonExistentUser() {
	body := strings.NewReader(`{
		"id": 30001,
		"username": "Surgeous",
		"password_hash": "abcxyz"}`)
	// Create a fake HTTP request
	req, _ := http.NewRequest("PATCH", "/user/update", body)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"no matched document found for update"}`, w.Body.String())
}

func (suite *UserTestSuite) TestUpdateSuccessfulUser() {
	body := strings.NewReader(`{
		"id": 20000,
		"username": "SurgeousJP",
		"password_hash": "abcxyz",
		"favorite_list": [2,3,5,6,11],
		"recommendation_list": [2,3,5,6,11],
		"watch_list": [2,3,5,6,11],
		"picture_profile": "" }`)

	req, _ := http.NewRequest("PATCH", "/user/update", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	assert.Equal(suite.T(), `{"message":"Successful"}`, w.Body.String())
}

func (suite *UserTestSuite) TestChangePasswordUnbindJSONUser() {
	body := strings.NewReader(`{UnbindedJSON}`)

	req, _ := http.NewRequest("PATCH", "/user/change_password", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"invalid character 'U' looking for beginning of object key string"}`, w.Body.String())
}

func (suite *UserTestSuite) TestChangePasswordWrongStructureKeyword() {
	body := strings.NewReader(`{
		"old_password": "12345",
		"new_password": "12345"
	}`)

	req, _ := http.NewRequest("PATCH", "/user/change_password", body)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	assert.Equal(suite.T(), `{"message":"Key: 'ResetPassword.ID' Error:Field validation for 'ID' failed on the 'required' tag"}`, w.Body.String())
}

func (suite *UserTestSuite) TestChangePasswordNonExistentUser() {
	body := strings.NewReader(`{
		"id": 30000,
		"old_password": "abcxyz",
		"new_password": "abcxyz"}`)
	// Create a fake HTTP request
	req, _ := http.NewRequest("PATCH", "/user/change_password", body)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"mongo: no documents in result"}`, w.Body.String())

}

func (suite *UserTestSuite) TestChangePasswordWrongOldPassword() {
	body := strings.NewReader(`{
		"id": 10001,
		"old_password": "1234",
		"new_password": "abcxyz"}`)
	// Create a fake HTTP request
	req, _ := http.NewRequest("PATCH", "/user/change_password", body)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"wrong password"}`, w.Body.String())
}

func (suite *UserTestSuite) TestDeleteInvalidUserIdUser() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("DELETE", "/user/delete/fsaf", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"Invalid movie id"}`, w.Body.String())
}

func (suite *UserTestSuite) TestDeleteNonExistentUser() {
	// Create a fake HTTP request
	req, _ := http.NewRequest("DELETE", "/user/delete/30000", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	suite.r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadGateway, w.Code)

	// Check the response body
	assert.Equal(suite.T(), `{"message":"no matched document found for delete"}`, w.Body.String())
}

func (suite *UserTestSuite) TestRegisterUserRouteSuccessfully() {
	// Create a new gin router
	router := gin.New()

	// Create an instance of your CastController
	cc := suite.userController // replace CastController with your actual type

	// Register the cast routes
	cc.RegisterUserRoute(router.Group("/v2"))

	// Perform a GET request to "/v2/cast/get/1"
	w := performUserRequest(router, "GET", "/v2/user/get/10001")

	// Assert the response status code
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// You can add more assertions for other routes if needed
}

func performUserRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func (suite *UserTestSuite) TearDownSuite() {
	req, _ := http.NewRequest("DELETE", "/user/delete/20000", nil)

	w := httptest.NewRecorder()

	suite.r.ServeHTTP(w, req)
}

func (suite *UserTestSuite) setupUserRouter() *gin.Engine {
	r := gin.New()
	r.GET("/user/get/:id", suite.userController.GetUser)
	r.POST("/user/create", suite.userController.CreateUser)
	r.PATCH("/user/update", suite.userController.UpdateUser)
	r.PATCH("/user/change_password", suite.userController.ChangePassword)
	r.DELETE("/user/delete/:id", suite.userController.DeleteUser)
	return r
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
