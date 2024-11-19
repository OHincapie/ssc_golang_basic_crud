package main

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math/rand"
	"net/http"
	"time"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type UserCRUD struct {
	userCollection *mongo.Collection
}

func (service *UserCRUD) GetUsers(c *gin.Context) {
	var users []*User
	ctx := c.Request.Context()
	filter := bson.D{{}}
	cursor, err := service.userCollection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	for cursor.Next(ctx) {
		var user User
		err := cursor.Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err,
			})
			return
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	cursor.Close(ctx)

	if len(users) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func (service *UserCRUD) CreateUser(c *gin.Context) {
	var user User
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	var userFetch *User
	query := bson.D{bson.E{Key: "email", Value: user.Email}}
	err := service.userCollection.FindOne(ctx, query).Decode(&userFetch)
	if err == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error user already exists",
		})
		return
	}
	user.Password = StringWithCharset(10, charset)
	_, insertErr := service.userCollection.InsertOne(ctx, user)
	if insertErr != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": insertErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  "success",
		"user":     user.Email,
		"password": user.Password,
	})
}

func (service *UserCRUD) Login(c *gin.Context) {
	var user User
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	var userFetch *User
	query := bson.D{bson.E{Key: "email", Value: user.Email}}
	err := service.userCollection.FindOne(ctx, query).Decode(&userFetch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error user not exists",
		})
		return
	}
	if userFetch.Email == user.Email && userFetch.Password == user.Password {
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{
		"message": "unauthorized",
	})
}

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
