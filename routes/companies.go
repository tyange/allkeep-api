package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tyange/white-shadow-api/models"
)

func createCompany(context *gin.Context) {
	var company models.Company
	err := context.ShouldBindBodyWithJSON(&company)

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request data."})
		return
	}

	userId := context.GetInt64("userId")
	company.UserID = userId

	err = company.Save()

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create company. Try again later."})
	}

	context.JSON(http.StatusCreated, gin.H{"message": "company created.", "company": company})
}
