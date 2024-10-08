package controllers

import (
	"gotest/models"
	"gotest/services/interfaces"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TopRatedMoviesController struct {
	TopRatedMoviesService interfaces.TopRatedMovieServices
	RatingMoviesService   interfaces.RatingService
	MovieService          interfaces.MovieService
}

func NewTopRatedMoviesController(TopRatedMoviesService interfaces.TopRatedMovieServices, RatingService interfaces.RatingService, MovieService interfaces.MovieService) TopRatedMoviesController {
	return TopRatedMoviesController{
		TopRatedMoviesService: TopRatedMoviesService,
		RatingMoviesService:   RatingService,
		MovieService:          MovieService,
	}
}

func (tc *TopRatedMoviesController) GetTopRatedMovies(ctx *gin.Context) {
	TopRatedMoviesUsers, _ := tc.TopRatedMoviesService.GetTopRatedMovies()

	ctx.JSON(http.StatusOK, TopRatedMoviesUsers)
}
func (tc *TopRatedMoviesController) CreateTopRatedMovies(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("movie_id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	movies, err := tc.MovieService.FindMovie(&id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	str := &id
	avg, err := tc.RatingMoviesService.GetAverageRating(str)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	topRatedMovies, err := tc.TopRatedMoviesService.FindTopRatedMovies(&id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if topRatedMovies != nil {
		topRatedMovies.VoteAvg = avg
		err = tc.TopRatedMoviesService.UpdateTopRatedMovies(topRatedMovies)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "TopRatedMovies updated successfully"})
		return
	}
	topRatedMoviesUser := &models.TopRatedMovies{
		MovieId:     id,
		VoteAvg:     avg,
		PosterPath:  movies.PosterPath,
		Title:       movies.Title,
		ReleaseDate: movies.ReleaseDate,
	}

	err = tc.TopRatedMoviesService.CreateTopRatedMovies(topRatedMoviesUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "TopRatedMovies created successfully"})
}
func (tc *TopRatedMoviesController) RegisterTopRatedMoviesRoute(rg *gin.RouterGroup) {
	TopRatedMoviesRoute := rg.Group("/topMovies")
	TopRatedMoviesRoute.GET("/get", tc.GetTopRatedMovies)
	TopRatedMoviesRoute.POST("/create/:movie_id", tc.CreateTopRatedMovies)
	// TopRatedMoviesRoute.GET("/get/:movie_id", tc.GetTopRatedMovies)
}
