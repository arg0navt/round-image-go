package user

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type FindItem struct {
	Title, Price, Location string
	Id                     int
}

func ParseUsers(w http.ResponseWriter, r *http.Request) {
	var (
		wg     sync.WaitGroup
		result []FindItem
	)
	for i := 100; i <= 112; i++ {
		wg.Add(1)
		go addItem(i, &wg, &result)
	}
	wg.Wait()
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
}

func GetParseUsers(url string, result *[]FindItem) {
	res, err := http.Get(url)
	if err == nil {
		doc, _ := goquery.NewDocumentFromReader(res.Body)
		doc.Find(".item").Each(func(i int, s *goquery.Selection) {
			id, okId := s.Attr("data-item-id")
			if okId == true {

				if s.Find(".item-description-title span").Text() != "" {
					newItem := FindItem{
						Title:    s.Find(".item-description-title span").Text(),
						Price:    s.Find(".about .price").Text(),
						Location: s.Find(".addres").Text(),
						Id:       strconv.Atoi(id),
					}
					*result = append(*result, newItem)
				}
			}
		})
		res.Body.Close()
	}
}

func findItemInResult(item int, result *[]FindItem) bool {
	for i := 0; i <= len(*result); i++ {
		if *result[i].id == item {
			return true
		}
	}
	return false
}

func addItem(i int, wg *sync.WaitGroup, result *[]FindItem) {
	GetParseUsers("https://www.avito.ru/rostov-na-donu/kvartiry?p="+strconv.Itoa(i), result)
	wg.Done()
}
