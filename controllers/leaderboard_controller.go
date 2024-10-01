package controllers

import (
	"gotest/services/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LeaderboardController struct {
	LeaderboardService interfaces.LeaderboardServices
}

func NewLeaderboardController(leaderboardService interfaces.LeaderboardServices) LeaderboardController {
	return LeaderboardController{
		LeaderboardService: leaderboardService,
	}
}

func (lc *LeaderboardController) GetLeaderboard(ctx *gin.Context) {
	leaderboardUsers, _ := lc.LeaderboardService.GetLeaderboard()

	ctx.JSON(http.StatusOK, leaderboardUsers)
}

func (lc *LeaderboardController) RegisterLeaderboardRoute(rg *gin.RouterGroup) {
	leaderboardRoute := rg.Group("/leaderboard")
	leaderboardRoute.GET("/get", lc.GetLeaderboard)
}