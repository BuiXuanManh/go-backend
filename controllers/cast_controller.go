package controllers

import (
	"gotest/models"
	"gotest/services/interfaces"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CastController struct {
	CastService interfaces.CastService
}

func NewCastController(castService interfaces.CastService) CastController {
	return CastController{
		CastService: castService,
	}
}

func (cc *CastController) CreateCast(ctx *gin.Context) {
	var cast models.Cast

	if err := ctx.ShouldBindJSON(&cast); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if cast.MovieId == 0 || len(cast.Cast) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "wrong input structure"})
		return
	}

	if err := cc.CastService.CreateCast(&cast); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successful"})
}

func (cc *CastController) GetCast(ctx *gin.Context) {
	movieId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	cast, err := cc.CastService.GetCast(&movieId)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, cast)
}

func (cc *CastController) UpdateCast(ctx *gin.Context) {
	var cast models.Cast

	if err := ctx.ShouldBindJSON(&cast); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if cast.MovieId == 0 || len(cast.Cast) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "wrong input structure"})
		return
	}

	if err := cc.CastService.UpdateCast(&cast); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successful"})
}

func (cc *CastController) DeleteCast(ctx *gin.Context) {
	movieId, err := strconv.Atoi(ctx.Param("id"))

	if err != nil || int64(movieId) <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid movie id"})
		return
	}
	
	err = cc.CastService.DeleteCast(&movieId)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successful"})
}

func (cc *CastController) RegisterCastRoute(rg *gin.RouterGroup) {
	castRoute := rg.Group("/cast")
	// The URI must be diffent structure from each other !
	castRoute.POST("/create", cc.CreateCast)

	castRoute.GET("/get/:id", cc.GetCast)

	castRoute.PATCH("/update", cc.UpdateCast)

	castRoute.DELETE("/delete/:id", cc.DeleteCast)
}
