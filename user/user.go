package user

import (
	"github.com/gin-gonic/gin"
	"log"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strconv"
)

type FindUser struct {
	Name, UrlProf, Avatar string
}

var users []FindUser

func GetParseUsers(url string, c chan []FindUser) {
	var arrayUser []FindUser
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	doc.Find(".shout-answer--item").Each(func(i int, s *goquery.Selection) {
		link, errLink := s.Find(".shout-author").Attr("href")
		img, errImg := s.Find(".shout-item--pic").Attr("src")
		if (errLink == true) && (errImg == true) {
			newUser := FindUser{
				Name: s.Find(".shout-author--name").Text(),
				UrlProf: link,
				Avatar: img,
			}
			arrayUser = append(arrayUser, newUser)
		}
	})
	c <- arrayUser
}

func addUsers(i int) {
	c := make(chan []FindUser)
	go GetParseUsers("https://kanobu.ru/shouts/" + strconv.Itoa(i), c)
	newUsers := <-c
	if newUsers != nil {
		users = append(users, newUsers...)
	}
}

func GetAll(c *gin.Context) {
	for i := 1; i < 20; i++ {
		go addUsers(i)
	}
	defer c.JSON(200, gin.H{
		"users": users,
	})
}