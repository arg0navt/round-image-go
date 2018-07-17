package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"fmt"
	"encoding/json"
)

type ResultIndex struct {
	Users []FindUser
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	makeWorker := Worker{Size: 20, Result: 0}
	var lookReady func(int)
	lookReady = func (i int) {
		fmt.Println(i)
		if i == makeWorker.Size {
			fmt.Println("piy piy")
			if err := json.NewEncoder(w).Encode(users); err != nil {
				panic(err)
			}
		}
	}
	for i := 1; i < makeWorker.Size; i++ {
		go addUsers(i, &makeWorker, lookReady)
	}
}

type Worker struct {
	Size int
	Result int
}

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

func addUsers(i int, worker *Worker, callback func(int)) {
	c := make(chan []FindUser)
	go GetParseUsers("https://kanobu.ru/shouts/" + strconv.Itoa(i), c)
	newUsers := <-c
	if newUsers != nil {
		users = append(users, newUsers...)
		worker.Result = worker.Result + 1
		callback(worker.Result)
	}
}