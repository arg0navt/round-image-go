package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"./user"
)

func startPage(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func main() {
	session, err := mgo.Dial("mongodb+srv://argonavt:korolego33@cluster0-vmmdn.mongodb.net/Cluster0?retryWrites=true")
    if err != nil {
        panic(err)
    }
    defer session.Close()
	route := gin.Default()
	v1 := route.Group("/api/v1/")
	route.GET("/", startPage)
	{
		v1.GET("user", user.GetUser)
	}
	route.Run()
}