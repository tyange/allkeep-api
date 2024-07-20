package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tyange/white-shadow-api/models"
)

type WorkStartBody struct {
	StartAt *time.Time `json:"start_at"`
	DoneAt  *time.Time `json:"done_at"`
}

type WorkPauseBody struct {
	PauseAt *time.Time `json:"pause_at"`
}

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
	work.IsPause = false

	err = work.Save()

	if err != nil {
		switch err.(type) {
		case *models.DuplicateCompanyIDError:
			fmt.Println(err)
			context.JSON(http.StatusConflict, gin.H{"message": "Company duplicated."})
			return
		default:
			fmt.Println(err)
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create work. Try again later."})
			return
		}
	}

	context.JSON(http.StatusCreated, gin.H{"message": "work created.", "work": work})
}

func updateWork(context *gin.Context) {
	workId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse work id."})
		return
	}

	userId := context.GetInt64("userId")
	work, err := models.GetWorkById(&workId)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch the work. Try again later."})
		return
	}

	if work.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized to update work."})
		return
	}

	var updatedWork models.Work
	err = context.ShouldBindBodyWithJSON(&updatedWork)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request data"})
		return
	}

	updatedWork.ID = workId
	err = updatedWork.Update()
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update the work. Try again later."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Work updated successfully!"})
}

func workStart(context *gin.Context) {
	workId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse work id."})
		return
	}

	userId := context.GetInt64("userId")
	work, err := models.GetWorkById(&workId)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch the work. Try again later."})
		return
	}

	if work.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized to update work."})
		return
	}

	var workStartBody WorkStartBody
	err = context.ShouldBindBodyWithJSON(&workStartBody)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request data"})
		return
	}

	if workStartBody.StartAt == nil || workStartBody.DoneAt == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "start_at and done_at are required"})
		return
	}

	err = models.UpdateWorkForStart(&workId, workStartBody.StartAt, workStartBody.DoneAt)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update the work. Try again later."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Work updated successfully!"})
}

func workPause(context *gin.Context) {
	workId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse work id."})
		return
	}

	userId := context.GetInt64("userId")
	work, err := models.GetWorkById(&workId)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch the work. Try again later."})
		return
	}

	if work.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized to update work."})
		return
	}

	var workPauseBody WorkPauseBody
	err = context.ShouldBindBodyWithJSON(&workPauseBody)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request data"})
		return
	}

	if workPauseBody.PauseAt == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "pause_at are required"})
		return
	}

	err = models.UpdateWorkForPause(&workId, workPauseBody.PauseAt)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update the work. Try again later."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Work updated successfully!"})
}
