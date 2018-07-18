package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"encoding/json"
	"sync"
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
	var (
		wg sync.WaitGroup
		result []FindUser
	)
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go addUsers(i, &wg, &result)
	}
	wg.Wait()
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
}

type FindUser struct {
	Name, UrlProf, Avatar string
}

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

func addUsers(i int, wg *sync.WaitGroup, result *[]FindUser) {
	c := make(chan []FindUser)
	go GetParseUsers("https://kanobu.ru/shouts/" + strconv.Itoa(i), c)
	newUsers := <-c
	if newUsers != nil {
		*result = append(*result, newUsers...)
	}
	wg.Done()
}