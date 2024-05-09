package controllers

import (
	"context"
	"net/http"
	// "strconv"
	"time"

	"github.com/Stingknight/expense-tracker/intializers"
	"github.com/Stingknight/expense-tracker/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddExpense()gin.HandlerFunc{
	return func(ctx *gin.Context){
		context,cancel := context.WithTimeout(context.Background(),10*time.Second)

		defer cancel()

		userDetails,exists := ctx.Get("user")
		if !exists{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":"Unautorized"})
			return
		}

		var userData = userDetails.(models.User)

		var expense models.ExpenseTrackModel

		if err := ctx.BindJSON(&expense);err!=nil{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":err})
			return
		}
		if err := validate.Struct(&expense);err!=nil{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":err})
			return
		}
		expense.UserExpense = userData.Id

		expense.Id  = primitive.NewObjectID()
		expense.CreateAt,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		expense.UpdatedAt,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))

		insertedResult,err := intializers.DB.Collection("expense").InsertOne(context,&expense)
		if err!=nil{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":err})
			return
		} 
		ctx.IndentedJSON(http.StatusOK,insertedResult)
	}	
}

func GetMyAllExpense()gin.HandlerFunc{
	return func(ctx *gin.Context){
		context,cancel := context.WithTimeout(context.Background(),10*time.Second)

		defer cancel()

		userDetails,exists := ctx.Get("user")
		if !exists{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":"Unautorized"})
			return
		}

		var userData = userDetails.(models.User)

		// var projectStage bson.D = bson.D{{Key: "$project",Value: bson.D{{Key: "username",Value: 1},{Key: "_id",Value: 1}}}}

		// var lookUpStage bson.D = bson.D{{Key: "$lookup",Value: bson.D{{Key: "from",Value: "user"},{Key: "foreignField",Value: "_id"},{Key: "localField",Value: "user_expense"},{Key: "pipeline",Value: []interface{}{projectStage}},{Key: "as",Value: "user_expense"}}}}

		// var unwindStage bson.D = bson.D{{Key: "$unwind",Value: bson.D{{Key: "path",Value: "$user_expense"}}}}

		var matchStage bson.D = bson.D{{Key: "$match",Value: bson.D{{Key: "user_expense",Value: userData.Id}}}}
		
		var groupStage bson.D = bson.D{{Key: "$group",Value: bson.D{{Key: "_id",Value: "$user_expense._id"},{Key: "user_expense_list",Value: bson.D{{Key: "$push",Value: "$$ROOT"}}},{Key: "totalSpent",Value: bson.D{{Key: "$sum",Value: "$expense_price"}}}}}}

		cursor,err := intializers.DB.Collection("expense").Aggregate(context,mongo.Pipeline{matchStage,groupStage})
		if err!=nil{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":err})
			return
		} 
		
		var expenseDetails []bson.M

		if err := cursor.All(context,&expenseDetails);err!=nil{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":err})
			return
		}
		ctx.IndentedJSON(http.StatusOK,expenseDetails)
	}


}
func GetMyExpenseType()gin.HandlerFunc{
	return func(ctx *gin.Context){
		context,cancel := context.WithTimeout(context.Background(),10*time.Second)

		defer cancel()

		userDetails,exists := ctx.Get("user")
		if !exists{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":"Unautorized"})
			return
		}

		var userData = userDetails.(models.User)

		// var projectStage bson.D = bson.D{{Key: "$project",Value: bson.D{{Key: "username",Value: 1},{Key: "_id",Value: 1}}}}

		// var lookUpStage bson.D = bson.D{{Key: "$lookup",Value: bson.D{{Key: "from",Value: "user"},{Key: "foreignField",Value: "_id"},{Key: "localField",Value: "user_expense"},{Key: "pipeline",Value: []interface{}{projectStage}},{Key: "as",Value: "user_expense"}}}}

		// var unwindStage bson.D = bson.D{{Key: "$unwind",Value: bson.D{{Key: "path",Value: "$user_expense"}}}}

		var matchStage bson.D = bson.D{{Key: "$match",Value: bson.D{{Key: "user_expense",Value: userData.Id}}}}
		
		var groupStage bson.D = bson.D{{Key: "$group",Value: bson.D{{Key: "_id",Value: "$expense_type"},{Key: "user_expense_list",Value: bson.D{{Key: "$push",Value: "$$ROOT"}}},{Key: "totalSpent",Value: bson.D{{Key: "$sum",Value: "$expense_price"}}}}}}

		var addFields bson.D = bson.D{{Key: "$addFields",Value: bson.D{{Key: "totalEntry",Value: bson.D{{Key: "$size",Value: "$user_expense_list"}}}}}}

		cursor,err := intializers.DB.Collection("expense").Aggregate(context,mongo.Pipeline{matchStage,groupStage,addFields})
		if err!=nil{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":err})
			return
		} 
		
		var expenseDetails []bson.M

		if err := cursor.All(context,&expenseDetails);err!=nil{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":err})
			return
		}
		ctx.IndentedJSON(http.StatusOK,expenseDetails)
	}


}



func UpdateExpense()gin.HandlerFunc{
	return func(ctx *gin.Context){
		context,cancel := context.WithTimeout(context.Background(),10*time.Second)

		defer cancel()

		userDetails,exists := ctx.Get("user")
		if !exists{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":"Unautorized"})
			return
		}

		var userData = userDetails.(models.User)

		var expense_id string = ctx.Param("expense_id")

		expenseObj_id,_ := primitive.ObjectIDFromHex(expense_id)

		var expense models.ExpenseTrackModel

		if err := ctx.BindJSON(&expense);err!=nil{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":err})
			return
		}

		var existExpense models.ExpenseTrackModel

		if err := intializers.DB.Collection("expense").FindOne(context,bson.M{"_id":expenseObj_id}).Decode(&existExpense);err!=nil{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":err})
			return
		}

		if existExpense.UserExpense!=userData.Id{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":"you cant update the expense"})
			return
		}

		var updateObj bson.D

		if expense.ExpenseType!=nil{
			updateObj = append(updateObj,bson.E{Key: "expense_type",Value: *expense.ExpenseType})
		}

		if expense.ExpensePrice!=nil{
			updateObj = append(updateObj,bson.E{Key: "expense_price",Value: *expense.ExpensePrice})
		}



		var filter bson.M = bson.M{"_id":expenseObj_id}

		var upsert bool = true

		options := options.UpdateOptions{
			Upsert: &upsert,
		}
		
		updateResult,err := intializers.DB.Collection("expense").UpdateOne(context,filter,bson.M{"$set":updateObj},&options)
		if err!=nil{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":err})
			return
		}
		ctx.IndentedJSON(http.StatusOK,updateResult) 
	}
}

func DeleteExpense()gin.HandlerFunc{
	return func(ctx *gin.Context){
		context,cancel := context.WithTimeout(context.Background(),10*time.Second)

		defer cancel()

		userDetails,exists := ctx.Get("user")
		if !exists{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":"Unautorized"})
			return
		}

		var userData = userDetails.(models.User)

		var expense_id string = ctx.Param("expense_id")

		expenseObj_id,_ := primitive.ObjectIDFromHex(expense_id)

		var existExpense models.ExpenseTrackModel

		if err := intializers.DB.Collection("expense").FindOne(context,bson.M{"_id":expenseObj_id}).Decode(&existExpense);err!=nil{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":err})
			return
		}

		if existExpense.UserExpense!=userData.Id{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":"you can delete this expense"})
			return
		}

		deleteResult,err := intializers.DB.Collection("expense").DeleteOne(context,bson.M{"_id":expenseObj_id})
		if err!=nil{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":err})
			return
		}
		ctx.IndentedJSON(http.StatusOK,deleteResult) 


	}


}

func GetMyMonthExpense()gin.HandlerFunc{

	return func(ctx *gin.Context){
		context,cancel := context.WithTimeout(context.Background(),10*time.Second)

		defer cancel()

		userDetails,exists := ctx.Get("user")
		if !exists{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":"Unautorized"})
			return
		}
		var userData = userDetails.(models.User)


		var matchStage bson.D = bson.D{{Key: "$match",Value: bson.D{{Key: "user_expense",Value: userData.Id}}}}
		
		var groupStage bson.D = bson.D{{Key: "$group",Value: bson.D{{Key: "_id",Value: bson.D{{Key: "$month",Value: "$created_at"}}},{Key: "totalsum",Value: bson.D{{Key: "$sum",Value: "$expense_price"}}}}}}

		cursor,err := intializers.DB.Collection("expense").Aggregate(context,mongo.Pipeline{matchStage,groupStage})
		if err!=nil{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":err})
			return
		}

		var expenseDetails []bson.M

		if err := cursor.All(context,&expenseDetails);err!=nil{
			ctx.IndentedJSON(http.StatusBadRequest,gin.H{"error":err})
			return
		}
	
		ctx.IndentedJSON(http.StatusOK,expenseDetails)
	}
}