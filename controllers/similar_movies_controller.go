package controllers

import (
	"gotest/services/interfaces"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SimilarMoviesController struct {
	SimilarMoviesService interfaces.SimilarMoviesServices
}

func NewSimilarMoviesController(similarMovieService interfaces.SimilarMoviesServices) SimilarMoviesController {
	return SimilarMoviesController{
		SimilarMoviesService: similarMovieService,
	}
}

func (lc *SimilarMoviesController) GetSimilarMovies(ctx *gin.Context) {
	movieId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	
	similarMovie, err := lc.SimilarMoviesService.GetSimilarMovies(&movieId)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, similarMovie)
}

func (lc *SimilarMoviesController) RegisterSimilarMoviesRoute(rg *gin.RouterGroup) {
	similarMovieRoute := rg.Group("/similarMovie")
	similarMovieRoute.GET("/get/:id", lc.GetSimilarMovies)
}