package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tyange/triplework-backend/models"
)

func workStart(context *gin.Context) {
	var work models.Work
	err := context.ShouldBindBodyWithJSON(&work)

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request data"})
		return
	}

	userId := 1
	work.UserID = int64(userId)

	err = work.Start()

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create work. Try again later."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "work created", "work": work})
}
