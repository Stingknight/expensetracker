package main

import (
	"github.com/Stingknight/expense-tracker/intializers"
	"github.com/Stingknight/expense-tracker/routes"
	"github.com/gin-gonic/gin"
	"log"
	"github.com/joho/godotenv"
)


func main(){
	
	 // Find .env file
	err := godotenv.Load(".env")
	if err != nil{
	  log.Fatalf("Error loading .env file: %s", err)
	}
	
	intializers.DatabaseIntialize()

	router := gin.Default()

	routes.UserRoutes(router)
	routes.ExpenseRoutes(router)

	router.Run("localhost:10000")
}