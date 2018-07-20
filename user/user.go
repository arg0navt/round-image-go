package user

import (
	"net/http"
	"sync"
	"strconv"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
)

type FindUser struct {
	Name, UrlProf, Avatar string
}

type SyntaxError struct {
	msg string // error description
}

const maxStack = 200

func ParseUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	start, errStart := strconv.Atoi(query["start"][0])
	end, errEnd := strconv.Atoi(query["end"][0])
	if errStart == nil && errEnd == nil {
		var (
			wg sync.WaitGroup
			result []FindUser
		)
		for i := start; i <= end; i++ {
			wg.Add(1)
			go addUsers(i, &wg, &result)
		}
		wg.Wait()
		if err := json.NewEncoder(w).Encode(result); err != nil {
			panic(err)
		}
	} else {
		if errEnd != nil {
			http.Error(w, errEnd.Error(), 500)
		}
		if errStart != nil {
			http.Error(w, errStart.Error(), 500)
		}
	}
}


func GetParseUsers(url string, result *[]FindUser) {
	res, err := http.Get(url)
	defer res.Body.Close()
	if err == nil {
		doc, _ := goquery.NewDocumentFromReader(res.Body)
		doc.Find(".shout-answer--item").Each(func(i int, s *goquery.Selection) {
			link, errLink := s.Find(".shout-author").Attr("href")
			img, errImg := s.Find(".shout-item--pic").Attr("src")
			if (errLink == true) && (errImg == true) {
				newUser := FindUser{
					Name: s.Find(".shout-author--name").Text(),
					UrlProf: link,
					Avatar: img,
				}
				if catchReiteration(result, newUser) != false {
					*result = append(*result, newUser)
				}
			}
		})
	}
}

func catchReiteration(arrayUser *[]FindUser, user FindUser) bool {
	for _, v := range *arrayUser {
		if v.Name == user.Name {
			return false
		}
	}
	return true
}

func addUsers(i int, wg *sync.WaitGroup, result *[]FindUser) {
	GetParseUsers("https://kanobu.ru/shouts/" + strconv.Itoa(i), result)
	wg.Done()
}