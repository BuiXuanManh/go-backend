package auth

import (
	"log"
	"gotest/controllers"
	"gotest/helper"
	"gotest/models"

	// "log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	// "github.com/joho/godotenv"
	jwt "github.com/appleboy/gin-jwt/v2"
)

var (
	identityKey = "userId"
	jwtSecret   []byte
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func NewJWTAuthMiddleware(userController *controllers.UserController) *jwt.GinJWTMiddleware {
	jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
	// Define the middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "validation zone",
		Key:         jwtSecret,
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					identityKey: v.UserId,
					// Add other claims as needed
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.User{
				UserId: int(claims[identityKey].(float64)),
				// Retrieve other claims as needed
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginVals.Username
			password := loginVals.Password

			user, err := userController.UserService.GetUserFromUsername(&username)

			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			if username == user.Username && helper.CheckPassword(user.PasswordHash, password) {

				return &models.User{
					UserId:         user.UserId,
					Username:       user.Username,
					PictureProfile: user.PictureProfile,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			user, _ := c.Get(identityKey)

			if v, ok := data.(*models.User); ok && v.UserId == user.(*models.User).UserId {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup: "header: Authorization, query: token, cookie: jwt",

		TokenHeadName: "Bearer",

		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}
	return authMiddleware
}
