package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func GetAll(c *gin.Context) {
	fmt.Println(c)
	c.JSON(200, gin.H{
		"user": 0,
	})
}