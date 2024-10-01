package controllers

import (
	"gotest/services/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TopRatedMoviesController struct {
	TopRatedMoviesService interfaces.TopRatedMovieServices
}

func NewTopRatedMoviesController(TopRatedMoviesService interfaces.TopRatedMovieServices) TopRatedMoviesController {
	return TopRatedMoviesController{
		TopRatedMoviesService: TopRatedMoviesService,
	}
}

func (tc *TopRatedMoviesController) GetTopRatedMovies(ctx *gin.Context) {
	TopRatedMoviesUsers, _ := tc.TopRatedMoviesService.GetTopRatedMovies()

	ctx.JSON(http.StatusOK, TopRatedMoviesUsers)
}

func (tc *TopRatedMoviesController) RegisterTopRatedMoviesRoute(rg *gin.RouterGroup) {
	TopRatedMoviesRoute := rg.Group("/topMovies")
	TopRatedMoviesRoute.GET("/get", tc.GetTopRatedMovies)
}