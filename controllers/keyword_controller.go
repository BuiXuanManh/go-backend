package controllers

import (
	"gotest/models"
	"gotest/services/interfaces"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type KeywordController struct {
	KeywordService interfaces.KeywordService
}

func NewKeywordController(keywordService interfaces.KeywordService) KeywordController {
	return KeywordController{
		KeywordService: keywordService,
	}
}

func (kc *KeywordController) CreateKeyword(ctx *gin.Context) {
	var keyword models.Keyword
	if err := ctx.ShouldBindJSON(&keyword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if keyword.MovieId == 0 || len(keyword.KeywordList) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "wrong input structure"})
		return
	}

	if err := kc.KeywordService.CreateKeyword(&keyword); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successful"})

}

func (kc *KeywordController) GetKeyword(ctx *gin.Context) {
	movieId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	keyword, err := kc.KeywordService.GetKeyword(&movieId)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, keyword)
}

func (kc *KeywordController) UpdateKeyword(ctx *gin.Context) {
	var keyword models.Keyword

	if err := ctx.ShouldBindJSON(&keyword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if keyword.MovieId == 0 || len(keyword.KeywordList) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "wrong input structure"})
		return
	}

	if err := kc.KeywordService.UpdateKeyword(&keyword); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successful"})
}

func (kc *KeywordController) DeleteKeyword(ctx *gin.Context) {
	movieId, err := strconv.Atoi(ctx.Param("id"))

	if err != nil || int64(movieId) <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid movie id"})
		return
	}
	err = kc.KeywordService.DeleteKeyword(&movieId)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successful"})
}

func (kc *KeywordController) RegisterKeywordRoute(rg *gin.RouterGroup) {
	keywordRoute := rg.Group("/keyword")
	// The URI must be diffent structure from each other !
	keywordRoute.POST("/create", kc.CreateKeyword)

	keywordRoute.GET("/get/:id", kc.GetKeyword)

	keywordRoute.PATCH("/update", kc.UpdateKeyword)

	keywordRoute.DELETE("/delete/:id", kc.DeleteKeyword)
}
