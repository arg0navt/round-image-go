package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type User struct {
    Email string
}

func GetAll(c *gin.Context) {
	fmt.Println(c)
	c.JSON(200, gin.H{
		"user": 0,
	})
}