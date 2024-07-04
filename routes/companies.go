package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tyange/white-shadow-api/models"
)

func getCompaniesByUserId(context *gin.Context) {
	userId := context.GetInt64("userId")

	pageSize, err := strconv.ParseInt(context.Query("pageSize"), 10, 64)
	if err != nil || pageSize <= 0 {
		pageSize = 4
	}

	pageNum, err := strconv.ParseInt(context.Query(("pageNum")), 10, 64)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	companies, err := models.GetAllCompanyByUserId(&userId, &pageSize, &pageNum)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch companies. Try again later."})
	}

	context.JSON(http.StatusOK, gin.H{"message": "get all companies by user.", "companies": companies})
}

func createCompany(context *gin.Context) {
	var company models.Company
	err := context.ShouldBindBodyWithJSON(&company)

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request data."})
		return
	}

	// auth middleware에서 만든 "userId" 데이터를 사용.
	userId := context.GetInt64("userId")
	company.UserID = userId

	err = company.Save()

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create company. Try again later."})
	}

	context.JSON(http.StatusCreated, gin.H{"message": "company created.", "company": company})
}

func updateCompany(context *gin.Context) {
	companyId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse company id."})
		return
	}

	userId := context.GetInt64("userId")
	company, err := models.GetCompanyById(&companyId)

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch the company. Try again later."})
		return
	}

	if company.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized to update company."})
		return
	}

	var updatedCompany models.Company
	err = context.ShouldBindBodyWithJSON(&updatedCompany)

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request data"})
		return
	}

	updatedCompany.ID = companyId
	err = updatedCompany.Update()

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update the event. Try again later."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Company updated successfully!"})
}
