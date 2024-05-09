package routes

import (
	"github.com/Stingknight/expense-tracker/controllers"
	"github.com/Stingknight/expense-tracker/middleware"
	"github.com/gin-gonic/gin"
)


func ExpenseRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.POST("/add-expense",middleware.ValidateToken(),controllers.AddExpense())
	incomingRoutes.GET("/get-expense",middleware.ValidateToken(),controllers.GetMyAllExpense())
	incomingRoutes.GET("/get-expense-type",middleware.ValidateToken(),controllers.GetMyExpenseType())
	incomingRoutes.PUT("/update-expense/:expense_id",middleware.ValidateToken(),controllers.UpdateExpense())
	incomingRoutes.DELETE("/delete-expense/:expense_id",middleware.ValidateToken(),controllers.DeleteExpense())
	incomingRoutes.GET("/get-expense-type-month/:month",middleware.ValidateToken(),controllers.GetMyMonthExpense())
}