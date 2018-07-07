package user

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"../mongo"
)

//type CollectionUser struct {
//	db *mgo.Collection
//}
//
//type User struct {
//	Id bson.ObjectId `bson:"_id"`
//	Email string `bson:"email"`
//}
//
//func GetCollectionUser(db *mgo.Database) {
//	collection = &CollectionUser{}
//}

func GetAll(c *gin.Context) {
	//usersList := []User{}
	//db.Find(bson.M{}).All(&usersList)
	fmt.Println(mongo.ConnectDb{})
	c.JSON(200, gin.H{
		"users": 0,
	})
}