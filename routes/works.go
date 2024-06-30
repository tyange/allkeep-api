package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tyange/white-shadow-api/models"
)

func getWorksByUserId(context *gin.Context) {
	userId := context.GetInt64("userId")

	works, err := models.GetAllWorksByUserId(&userId)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch works. Try again later."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "get all works by user.", "works": works})
}

func createWork(context *gin.Context) {
	var work models.Work
	err := context.ShouldBindBodyWithJSON(&work)

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request data"})
		return
	}

	// auth middleware에서 만든 "userId" 데이터를 사용.
	userId := context.GetInt64("userId")
	work.UserID = int64(userId)

	err = work.Save()

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create work. Try again later."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "work created.", "work": work})
}
