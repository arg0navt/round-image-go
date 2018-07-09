package main

import (
	"github.com/gin-gonic/gin"
	"./user"
	"./mongo"
)



func startPage(c *gin.Context) {
	c.JSON(200, gin.H{
		"connect": 1,
	})
}

func main() {
	session, err := mongo.NewSession("mongodb://localhost")
  	if(err != nil) {
		panic(err)
	}
  	defer session.Close()
	route := gin.Default()
	v1 := route.Group("/users")
	route.GET("/", startPage)
	{
		v1.GET("/", user.GetAll)
	}
	route.Run()
}