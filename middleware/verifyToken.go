package middleware

import (
	"context"
	"net/http"
	"github.com/Stingknight/expense-tracker/helpers"
	"github.com/Stingknight/expense-tracker/intializers"
	"github.com/Stingknight/expense-tracker/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func ValidateToken()gin.HandlerFunc{

	return func(ctx *gin.Context){
		tokenString,err := ctx.Cookie("Authorization")
		if err!=nil{
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims,err := helpers.TokenValidation(tokenString)
		if err!=nil{
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var UserId string = claims["Id"].(string)

		userObjectId,_ := primitive.ObjectIDFromHex(UserId)

		var user models.User

		if err := intializers.DB.Collection("user").FindOne(context.Background(),bson.M{"_id":userObjectId}).Decode(&user);err!=nil{
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("user",user)
		ctx.Next()
		
	}


}
