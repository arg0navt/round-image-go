package main

import (
	"github.com/gin-gonic/gin"
	"./user"
	"./mongo"
	"net/http"
	"log"
	"github.com/PuerkitoBio/goquery"
)

type FindUser struct {
	Name string
	UrlProf string
	Avatar string
}

func startPage(c *gin.Context) {
	res, err := http.Get("https://kanobu.ru/shouts/all")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	var arrayUser []FindUser
	doc.Find(".page-shouts .shout-item").Each(func(i int, s *goquery.Selection) {
		newUser := FindUser{Name: s.Find(".shout-author--name").Text()}
		arrayUser = append(arrayUser, newUser)
	})
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	c.JSON(200, gin.H{
		"connect": arrayUser,
	})
}

func main() {
	session, err := mongo.NewSession("mongodb://localhost")
  	if(err != nil) {
		panic(err)
	}
  	defer session.Close()
	route := gin.Default()
	v1 := route.Group("/user")
	route.GET("/", startPage)
	{
		v1.GET("/", user.GetAll)
	}
	route.Run()
}