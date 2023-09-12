package main

import (
	"context"
	"net/http"
	"time"
	"log"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
)
func setupDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	database = client.Database("newintern")
	collection = database.Collection("newintern")
}

type Interns struct {
	UserID   int    `json:"UserId" bson:"UserId"`
	Name string `json:"Name" bson:"Name"`
}

func GetIntern(c *gin.Context) {
	UserId := c.Param("userID")
	var intern Interns
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, bson.M{"ID": UserId}).Decode(&intern)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Intern not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	c.JSON(http.StatusOK, intern)

}

func UpdateIntern(c *gin.Context) {
	UserId := c.Param("UserID")
	var updatedintern Interns
	if err := c.ShouldBindJSON(&updatedintern); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	 
	update := bson.M{
		"$set": bson.M{
			"Name":  updatedintern.Name,
			"id": updatedintern.UserID,
			
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"Id": UserId}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	c.JSON(http.StatusOK, result)




	
}


func CreateIntern(c *gin.Context){
	var intern []Interns
	if err := c.ShouldBindJSON(&intern); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, intern)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	c.JSON(http.StatusOK, result)






}

func DeleteIntern(c *gin.Context) {
	UserId := c.Param("Id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := collection.DeleteOne(ctx, bson.M{"Id": UserId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func main() {
	router := gin.Default()

	router.GET("/api", GetIntern)
	router.POST("/create-intern", CreateIntern)
	router.PUT("/update-intern/:UserId", UpdateIntern)
	router.DELETE("/api/:UserId", DeleteIntern)

	router.Run("localhost:8080")

}
