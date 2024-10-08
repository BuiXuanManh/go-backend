package controllers

import (
	"gotest/models"
	"gotest/services/interfaces"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LeaderboardController struct {
	LeaderboardService interfaces.LeaderboardServices
	UserService        interfaces.UserService
	MovieService       interfaces.MovieService
	RatingService      interfaces.RatingService
}

func NewLeaderboardController(leaderboardService interfaces.LeaderboardServices, UserService interfaces.UserService, MovieService interfaces.MovieService, RatingService interfaces.RatingService) LeaderboardController {
	return LeaderboardController{
		LeaderboardService: leaderboardService,
		UserService:        UserService,
		MovieService:       MovieService,
		RatingService:      RatingService,
	}
}

func (lc *LeaderboardController) GetLeaderboard(ctx *gin.Context) {
	leaderboardUsers, _ := lc.LeaderboardService.GetLeaderboard()

	ctx.JSON(http.StatusOK, leaderboardUsers)
}
func (lc *LeaderboardController) CreateLeaderboard(ctx *gin.Context) {
	user_id, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// rate, err := strconv.Atoi(ctx.Param("rate"))
	user, err := lc.UserService.GetUser(&user_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rates, err := lc.RatingService.GetRatingOfUser(&user_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	leaderboardUserMovie, err := lc.LeaderboardService.FindLeaderboard(&user_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if leaderboardUserMovie != nil {
		err = lc.LeaderboardService.UpdateLeaderboard(&models.Leaderboard{
			MoviesRated: len(rates),
			// MovieID:        movie_id,
			UserId:         user_id,
			PictureProfile: user.PictureProfile,
			Username:       user.Username,
		}, &user_id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Leaderboard update successfully"})
		return
	}
	leaderboardUser := &models.Leaderboard{
		// MovieID:        movie_id,
		UserId:         user_id,
		MoviesRated:    1,
		Username:       user.Username,
		PictureProfile: user.PictureProfile,
	}
	err = lc.LeaderboardService.CreateLeaderboard(leaderboardUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Leaderboard created successfully"})
}

func (lc *LeaderboardController) RegisterLeaderboardRoute(rg *gin.RouterGroup) {
	leaderboardRoute := rg.Group("/leaderboard")
	leaderboardRoute.GET("/get", lc.GetLeaderboard)
	leaderboardRoute.POST("/create/:user_id", lc.CreateLeaderboard)
}
