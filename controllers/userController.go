package controllers

import (
	"context"
	
	"log"
	"net/http"
	"time"
	"github.com/Stingknight/expense-tracker/helpers"
	"github.com/Stingknight/expense-tracker/intializers"
	"github.com/Stingknight/expense-tracker/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate *validator.Validate = validator.New()

var password helpers.PasswordM

func Signup()gin.HandlerFunc{
	return func(ctx *gin.Context){
		context,cancel := context.WithTimeout(context.Background(),10*time.Second)

		defer  cancel()

		var user models.User

		err := ctx.BindJSON(&user)
		if err!=nil{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":err})
			return
		}

		if err := validate.Struct(&user);err!=nil{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":err})
			return
		}

		password = helpers.PasswordM{
			Password:*user.Password,
		}

		count,err := intializers.DB.Collection("user").CountDocuments(context,bson.M{"$or":[]bson.M{{"username":*user.Username},{"mobile":*user.Mobile}}})
		if err!=nil{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":err})
			return
		}
		
		if count > 0{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":"Username or mobile already exists"})
			return
		}

		hashedPassword,err := password.CreatePasswordHash()
		if err!=nil{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":err})
			return
		}
		
		user.Id = primitive.NewObjectID()
		user.CreateAt,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		user.UpdatedAt,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		user.Password = &hashedPassword
		
		insertedResult,err := intializers.DB.Collection("user").InsertOne(context,&user)
		if err!=nil{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":err})
			return
		}
		ctx.IndentedJSON(http.StatusOK,insertedResult)

	}
	
}

func Login()gin.HandlerFunc{
	return func(ctx *gin.Context){
		context,cancel := context.WithTimeout(context.Background(),10*time.Second)

		defer cancel()

		var user models.User

		if err := ctx.BindJSON(&user);err!=nil{
			log.Fatalf("error: %v", err)
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":err})
			return
		}

		var existUser models.User

		if err := intializers.DB.Collection("user").FindOne(context,bson.M{"username":*user.Username}).Decode(&existUser);err!=nil{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":"Username is wrong"})
			return
		}

		password =helpers.PasswordM{
			Password: *user.Password,
			HashedPassword: *existUser.Password,
		}

		_,err := password.ComparePassword()
		if err!=nil{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":"Password is wrong"})
			return
		}
		
		tokenString,err := helpers.GenerateToken(existUser.Id)
		if err!=nil{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":err})
			return
		} 

		ctx.SetSameSite(http.SameSiteLaxMode)
		ctx.SetCookie("Authorization", tokenString, 3600, "", "", false, true)
		ctx.IndentedJSON(http.StatusOK, gin.H{})
	}

}

func Validate()gin.HandlerFunc{
	return func(ctx *gin.Context){
		
		userDetails,exists := ctx.Get("user")
		if !exists{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error1":"Unautorized"})
			return
		}

		var userData = userDetails.(models.User)

		ctx.IndentedJSON(http.StatusOK,userData)

	}	
}