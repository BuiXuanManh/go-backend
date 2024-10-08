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
}

func NewLeaderboardController(leaderboardService interfaces.LeaderboardServices, UserService interfaces.UserService, MovieService interfaces.MovieService) LeaderboardController {
	return LeaderboardController{
		LeaderboardService: leaderboardService,
		UserService:        UserService,
		MovieService:       MovieService,
	}
}

func (lc *LeaderboardController) GetLeaderboard(ctx *gin.Context) {
	leaderboardUsers, _ := lc.LeaderboardService.GetLeaderboard()

	ctx.JSON(http.StatusOK, leaderboardUsers)
}
func (lc *LeaderboardController) CreateLeaderboard(ctx *gin.Context) {
	movie_id, err := strconv.Atoi(ctx.Param("movie_id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user_id, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rate, err := strconv.Atoi(ctx.Param("rate"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user, err := lc.UserService.GetUser(&user_id)
	// rate_id, err := strconv.Atoi(ctx.Param("rate_id"))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	leaderboardUserMovie, err := lc.LeaderboardService.FindLeaderboard(&movie_id, &user_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if leaderboardUserMovie != nil {
		leaderboardUserMovie.MoviesRated = rate
		err = lc.LeaderboardService.UpdateLeaderboard(leaderboardUserMovie, &movie_id, &user_id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Leaderboard update successfully"})
		return
	}
	leaderboardUser := &models.Leaderboard{
		MovieID:        movie_id,
		UserId:         user_id,
		MoviesRated:    rate,
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
func (lc *LeaderboardController) UpdateLeaderboard(ctx *gin.Context) {
	var leaderboardUser *models.Leaderboard
	ctx.BindJSON(&leaderboardUser)
	movie_id, err := strconv.Atoi(ctx.Param("movie_id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user_id, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = lc.LeaderboardService.UpdateLeaderboard(leaderboardUser, &movie_id, &user_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Leaderboard updated successfully"})
}
func (lc *LeaderboardController) RegisterLeaderboardRoute(rg *gin.RouterGroup) {
	leaderboardRoute := rg.Group("/leaderboard")
	leaderboardRoute.GET("/get", lc.GetLeaderboard)
	leaderboardRoute.POST("/create/:movie_id/:user_id/:rate", lc.CreateLeaderboard)
}
