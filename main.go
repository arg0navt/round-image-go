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
	session, err := mgo.Dial("mongodb://localhost/rimg")
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