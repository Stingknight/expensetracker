package intializers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)


var DB *mongo.Database 

func DatabaseIntialize(){
	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second)

	defer cancel()

	client,err := mongo.Connect(ctx,options.Client().ApplyURI(os.Getenv("DATABASE")))
	if err!=nil{
		log.Fatalf("error %v",err )
		return
	}

	DB = client.Database("expensetracker")

	fmt.Println("Connecting to database------------>")

	err = client.Ping(ctx,readpref.Primary())
	if err!=nil{
		log.Fatalf("error %v",err )
		return
	}

	fmt.Println("Connected to database------------>")
	
}

