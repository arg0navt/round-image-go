package user

import (
	"github.com/gin-gonic/gin"
	"../parseUsers"
)

func GetAll(c *gin.Context) {
	parseUsers.Start(10)
	c.JSON(200, gin.H{
		"users": "not",
	})
}