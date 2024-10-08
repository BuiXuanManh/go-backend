package main

import (
	"context"
	"gotest/controllers"
	auth "gotest/jwt-authenticate"
	"gotest/models"
	"gotest/services/implementations"
	"gotest/services/interfaces"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server *gin.Engine

	databaseName string

	movieCollection *mongo.Collection
	movieService    interfaces.MovieService
	movieController controllers.MovieController

	keywordCollection *mongo.Collection
	keywordService    interfaces.KeywordService
	keywordController controllers.KeywordController

	ratingCollection *mongo.Collection
	ratingService    interfaces.RatingService
	ratingController controllers.RatingController

	crewCollection *mongo.Collection
	crewService    interfaces.CrewService
	crewController controllers.CrewController

	castCollection *mongo.Collection
	castService    interfaces.CastService
	castController controllers.CastController

	userCollection *mongo.Collection
	userService    interfaces.UserService
	userController controllers.UserController

	movieDiscussionCollection *mongo.Collection
	movieDiscussionService    interfaces.MovieDiscussionServices
	movieDiscussionController controllers.MovieDiscussionController

	leaderboardCollection *mongo.Collection
	leaderboardService    interfaces.LeaderboardServices
	leaderboardController controllers.LeaderboardController

	topRatedMoviesCollection *mongo.Collection
	topRatedMoviesService    interfaces.TopRatedMovieServices
	topRatedMoviesController controllers.TopRatedMoviesController

	userReviewCollection *mongo.Collection
	userReviewService    interfaces.UserReviewService
	userReviewController controllers.UserReviewController

	similarMoviesCollection *mongo.Collection
	similarMoviesService    interfaces.SimilarMoviesServices
	similarMoviesController controllers.SimilarMoviesController

	ctx         context.Context
	mongoClient *mongo.Client
)

const (
	NUMBER_OF_SECONDS_IN_ONE_DAY = 86400
	DEFAULT_DATABASE_CODE        = 0
	DEFAULT_PORT                 = "9090"
)

func Init() {
	ctx = context.TODO()

	// You should load the DB_CONNECTION_STRING, SERVER_GROUP, PORT from your .env environment,
	// set it to fit your usage, you can check it from previous commits for more information

	// Load environment variables from the .env file (in local)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve the connection string from the environment
	connectionString := os.Getenv("DB_CONNECTION_STRING")

	mongoConn := options.
		Client().
		ApplyURI(connectionString)

	mongoClient, err = mongo.Connect(ctx, mongoConn)
	if err != nil {
		log.Fatal(err)
	}
	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	databaseName = os.Getenv("DB_NAME")

	movieCollection = mongoClient.Database(databaseName).Collection("movies")
	movieService = implementations.NewMovieService(movieCollection, ctx)
	movieController = controllers.NewMovieController(movieService)

	keywordCollection = mongoClient.Database(databaseName).Collection("keywords")
	keywordService = implementations.NewKeywordService(keywordCollection, ctx)
	keywordController = controllers.NewKeywordController(keywordService)

	ratingCollection = mongoClient.Database(databaseName).Collection("ratings")
	ratingService = implementations.NewRatingService(ratingCollection, ctx)
	ratingController = controllers.NewRatingController(ratingService)

	crewCollection = mongoClient.Database(databaseName).Collection("crews")
	crewService = implementations.NewCrewService(crewCollection, ctx)
	crewController = controllers.NewCrewController(crewService)

	castCollection = mongoClient.Database(databaseName).Collection("casts")
	castService = implementations.NewCastService(castCollection, ctx)
	castController = controllers.NewCastController(castService)

	userCollection = mongoClient.Database(databaseName).Collection("users")
	userService = implementations.NewUserService(userCollection, ctx)
	userController = controllers.NewUserController(userService)

	movieDiscussionCollection = mongoClient.Database(databaseName).Collection("discussion_movie")
	movieDiscussionService = implementations.NewMovieDiscussionServices(movieDiscussionCollection, ctx)
	movieDiscussionController = controllers.NewMovieDiscussionController(movieDiscussionService)

	leaderboardCollection = mongoClient.Database(databaseName).Collection("leaderboard")
	leaderboardService = implementations.NewLeaderboardService(leaderboardCollection, ctx)
	leaderboardController = controllers.NewLeaderboardController(leaderboardService, userService, movieService)

	topRatedMoviesCollection = mongoClient.Database(databaseName).Collection("top_rated_movies")
	topRatedMoviesService = implementations.NewTopRatedMoviesService(topRatedMoviesCollection, ctx)
	topRatedMoviesController = controllers.NewTopRatedMoviesController(topRatedMoviesService, ratingService, movieService)

	userReviewCollection = mongoClient.Database(databaseName).Collection("user_reviews")
	userReviewService = implementations.NewUserReviewService(userReviewCollection, ctx)
	userReviewController = controllers.NewUserReviewController(userReviewService)

	similarMoviesCollection = mongoClient.Database(databaseName).Collection("similar_movies")
	similarMoviesService = implementations.NewSimilarMoviesService(similarMoviesCollection, ctx)
	similarMoviesController = controllers.NewSimilarMoviesController(similarMoviesService)

	server = gin.Default()
	// Set up CORS (Cross-Origin Resource Sharing)
	config := cors.DefaultConfig()

	// Allow all origins
	config.AllowAllOrigins = true

	// Allow the "Authorization" header
	config.AddAllowHeaders("Authorization")

	server.Use(cors.New(config))

	server.SetTrustedProxies(nil)
}

func returnUser(c *gin.Context) {
	// claims := jwt.ExtractClaims(c)
	user, _ := c.Get("userId")
	c.JSON(200, gin.H{
		"userId": user.(*models.User).UserId,
	})
}

func main() {
	Init()

	authMiddleware := auth.NewJWTAuthMiddleware(&userController)
	defer mongoClient.Disconnect(ctx)

	serverGroup := os.Getenv("SERVER_GROUP")

	port := os.Getenv("PORT")
	if port == "" {
		port = DEFAULT_PORT
	}

	basePath := server.Group(serverGroup)

	// Register your routes
	basePath.POST("/login", authMiddleware.LoginHandler)

	basePath.GET("/refresh_token", authMiddleware.RefreshHandler)

	// Apply middleware only to the /currentUser route
	basePath.GET("/currentUser", authMiddleware.MiddlewareFunc(), returnUser)

	// Apply middleware to all routes under basePath, except for GET requests, and /v1/user/create
	basePath.Use(func(c *gin.Context) {
		if c.Request.Method != "GET" && c.FullPath() != "/v1/user/create" {
			authMiddleware.MiddlewareFunc()(c)
		}
	})

	// Other route registrations without the middleware
	userController.RegisterUserRoute(basePath)
	movieController.RegisterMovieRoute(basePath)
	keywordController.RegisterKeywordRoute(basePath)
	crewController.RegisterCrewRoute(basePath)
	castController.RegisterCastRoute(basePath)
	ratingController.RegisterRatingRoute(basePath)
	movieDiscussionController.RegisterMovieDiscussionRoute(basePath)
	leaderboardController.RegisterLeaderboardRoute(basePath)
	topRatedMoviesController.RegisterTopRatedMoviesRoute(basePath)
	userReviewController.RegisterUserReviewRoute(basePath)
	similarMoviesController.RegisterSimilarMoviesRoute(basePath)

	log.Fatal(server.Run(":" + port))
}
