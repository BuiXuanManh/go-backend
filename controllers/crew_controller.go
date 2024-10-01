package controllers

import (
	"gotest/models"
	"gotest/services/interfaces"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CrewController struct {
	CrewService interfaces.CrewService
}

func NewCrewController(crewService interfaces.CrewService) CrewController {
	return CrewController{
		CrewService: crewService,
	}
}

func (cc *CrewController) CreateCrew(ctx *gin.Context) {
	var crew models.Crew
	if err := ctx.ShouldBindJSON(&crew); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if crew.MovieId == 0 || len(crew.Crew) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "wrong input structure"})
		return
	}

	if err := cc.CrewService.CreateCrew(&crew); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successful"})
}

func (cc *CrewController) GetCrew(ctx *gin.Context) {
	movieId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	crew, err := cc.CrewService.GetCrew(&movieId)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, crew)
}

func (cc *CrewController) UpdateCrew(ctx *gin.Context) {
	var crew models.Crew

	if err := ctx.ShouldBindJSON(&crew); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if crew.MovieId == 0 || len(crew.Crew) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "wrong input structure"})
		return
	}

	if err := cc.CrewService.UpdateCrew(&crew); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successful"})
}

func (cc *CrewController) DeleteCrew(ctx *gin.Context) {
	movieId, err := strconv.Atoi(ctx.Param("id"))

	if err != nil || int64(movieId) <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid movie id"})
		return
	}
	err = cc.CrewService.DeleteCrew(&movieId)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successful"})
}

func (cc *CrewController) RegisterCrewRoute(rg *gin.RouterGroup) {
	crewRoute := rg.Group("/crew")
	// The URI must be diffent structure from each other !
	crewRoute.POST("/create", cc.CreateCrew)

	crewRoute.GET("/get/:id", cc.GetCrew)

	crewRoute.PATCH("/update", cc.UpdateCrew)

	crewRoute.DELETE("/delete/:id", cc.DeleteCrew)
}
