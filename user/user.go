package user

import (
	"github.com/gin-gonic/gin"
	"log"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"fmt"
)

type FindUser struct {
	Name, UrlProf, Avatar string
}

func GetAll(c *gin.Context) {
	res, err := http.Get("https://kanobu.ru/shouts/all")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	var arrayUser []FindUser
	doc.Find(".page-shouts .shout-item").Each(func(i int, s *goquery.Selection) {
		link, errLink := s.Find(".shout-author").Attr("href")
		img, errImg := s.Find(".shout-item--pic").Attr("src")
		if (errLink == true) && (errImg == true) {
			newUser := FindUser{Name: s.Find(".shout-author--name").Text(), UrlProf: link, Avatar: img}
			arrayUser = append(arrayUser, newUser)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	fmt.Println(arrayUser)
	c.JSON(200, gin.H{
		"users": arrayUser,
	})
}