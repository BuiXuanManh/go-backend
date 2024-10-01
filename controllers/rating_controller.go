package controllers

import (
	"fmt"
	"gotest/models"
	"gotest/services/interfaces"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RatingController struct {
	RatingService interfaces.RatingService
}

func NewRatingController(ratingService interfaces.RatingService) RatingController {
	return RatingController{
		RatingService: ratingService,
	}
}

func (rc *RatingController) CreateRating(ctx *gin.Context) {
	var rating models.Rating
	if err := ctx.ShouldBindJSON(&rating); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if rating.UserId == 0 || rating.MovieId == 0 || rating.Timestamp == 0 || rating.Rating == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "wrong input structure"})
		return
	}

	if err := rc.RatingService.CreateRating(&rating); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successful"})

}

func (rc *RatingController) UpdateRating(ctx *gin.Context) {
	var rating models.Rating

	if err := ctx.ShouldBindJSON(&rating); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if rating.UserId == 0 || rating.MovieId == 0 || rating.Timestamp == 0 || rating.Rating == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "wrong input structure"})
		return
	}

	if err := rc.RatingService.UpdateRating(&rating); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successful"})
}

func (rc *RatingController) DeleteRating(ctx *gin.Context) {
	movieId, err := strconv.Atoi(ctx.Param("movieId"))
	if err != nil || int64(movieId) <= 0 {
		fmt.Println(movieId)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid movie id"})
		return
	}

	userId, err := strconv.Atoi(ctx.Param("userId"))
	if err != nil || int64(userId) <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user id"})
		return
	}

	err = rc.RatingService.DeleteRating(&movieId, &userId)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successful"})
}

func (rc *RatingController) GetRatingOfMovie(ctx *gin.Context) {
	movieId, err := strconv.Atoi(ctx.Param("movie_id"))
	if err != nil || int64(movieId) <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid movie id"})
		return
	}

	ratings, err := rc.RatingService.GetRatingOfMovie(&movieId)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, ratings)
}

func (rc *RatingController) GetRatingOfUser(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil || int64(userId) <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid movie id"})
		return
	}

	ratings, err := rc.RatingService.GetRatingOfUser(&userId)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, ratings)
}

func (rc *RatingController) GetMovieRatingOfUser(ctx *gin.Context) {
	movieId, err := strconv.Atoi(ctx.Param("movie_id"))
	if err != nil || int64(movieId) <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid movie id"})
		return
	}

	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil || int64(userId) <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user id"})
		return
	}

	rating, err := rc.RatingService.GetMovieRatingOfUser(&movieId, &userId)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, rating)
}

func (rc *RatingController) RegisterRatingRoute(rg *gin.RouterGroup) {
	ratingRoute := rg.Group("/rating")
	// The URI must be diffent structure from each other !
	ratingRoute.POST("/create", rc.CreateRating)

	ratingRoute.GET("/get/movie/:movie_id", rc.GetRatingOfMovie)

	ratingRoute.GET("/get/user/:user_id", rc.GetRatingOfUser)

	ratingRoute.GET("/get/:user_id/:movie_id", rc.GetMovieRatingOfUser)

	ratingRoute.PATCH("/update", rc.UpdateRating)

	ratingRoute.DELETE("/delete/:userId/:movieId", rc.DeleteRating)
}
