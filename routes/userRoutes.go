package routes

import (
	"github.com/Stingknight/expense-tracker/controllers"
	"github.com/Stingknight/expense-tracker/middleware"
	"github.com/gin-gonic/gin"
)



func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.POST("/signup",controllers.Signup())
	incomingRoutes.POST("/login",controllers.Login())
	incomingRoutes.GET("/validate",middleware.ValidateToken(),controllers.Validate())
}
