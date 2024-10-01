package controllers

import (
	"gotest/models"
	"gotest/services/interfaces"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserReviewController struct {
	UserReviewService interfaces.UserReviewService
}

func NewUserReviewController(UserReviewService interfaces.UserReviewService) UserReviewController {
	return UserReviewController{
		UserReviewService: UserReviewService,
	}
}

func (uc *UserReviewController) CreateUserReview(ctx *gin.Context) {
	var UserReview models.UserReview

	if err := ctx.ShouldBindJSON(&UserReview); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := uc.UserReviewService.CreateUserReview(&UserReview); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successful"})
}

func (uc *UserReviewController) UpdateUserReview(ctx *gin.Context) {
	var UserReview models.UserReview

	if err := ctx.ShouldBindJSON(&UserReview); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := uc.UserReviewService.UpdateUserReview(&UserReview); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successful"})
}

func (uc *UserReviewController) DeleteUserReview(ctx *gin.Context) {
	movieId, err1 := strconv.Atoi(ctx.Param("movie_id"))
	userId, err2 := strconv.Atoi(ctx.Param("user_id"))
	if err1 != nil || err2 != nil || int64(movieId) <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid movie id"})
		return
	}
	
	err1 = uc.UserReviewService.DeleteUserReview(&userId, &movieId)

	if err1 != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err1.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successful"})
}

func (uc *UserReviewController) GetUserReviewByMovieId(ctx *gin.Context) {
	movieId, err := strconv.Atoi(ctx.Param("movie_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	userReviews, err := uc.UserReviewService.GetUserReviewByMovieId(&movieId)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, userReviews)
}

func (uc *UserReviewController) GetUserReviewByUserId(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	userReviews, err := uc.UserReviewService.GetUserReviewByUserId(&userId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, userReviews)
}

func (uc *UserReviewController) RegisterUserReviewRoute(rg *gin.RouterGroup) {
	UserReviewRoute := rg.Group("/userReview")
	// The URI must be diffent structure from each other !
	
	UserReviewRoute.GET("/get/user/:user_id", uc.GetUserReviewByUserId)

	UserReviewRoute.GET("/get/movie/:movie_id", uc.GetUserReviewByMovieId)
	
	UserReviewRoute.POST("/create", uc.CreateUserReview)

	UserReviewRoute.PATCH("/update", uc.UpdateUserReview)

	UserReviewRoute.DELETE("/delete/:user_id/:movie_id", uc.DeleteUserReview)
}
