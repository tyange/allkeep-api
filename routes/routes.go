package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tyange/white-shadow-api/middlewares"
)

func RegisterRoutes(server *gin.Engine) {
	// TODO: get work
	// TODO: get works
	// TODO: edit work
	// TODO: delete work, work에 대한 id를 storage에 저장하고 있다가 재설정을 누르면 해당 id의 work를 삭제.
	work := server.Group("/works")
	work.Use(middlewares.Authenticate)
	work.GET("/all", getWorksByUserId)
	work.POST("/create", createWork)

	// TODO: get user info
	// TODO: edit user info
	// TODO: delete user
	auth := server.Group("/auth")
	auth.POST("/signup", signup)
	auth.POST("/login", login)
	auth.POST("/google", googleLoginCallBack)

	// 일하는 곳은 복수일 수 있음.
	// TODO: get company
	// TODO: edit company
	// TODO: delete company
	company := server.Group("/companies")
	company.Use(middlewares.Authenticate)
	company.GET("/all", getCompaniesByUserId)
	company.GET("/all-at-once", getAllCompaniesByUserId)
	company.POST("/create", createCompany)
	company.PUT("/:id", updateCompany)
	company.DELETE("/:id", deleteCompany)
}
