package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type FindItem struct {
	Title, Price, Location string
	Id                     int
}

func ParseUsers(w http.ResponseWriter, r *http.Request) {
	var result []FindItem
	group := make(chan int)
	size := 200
	endItem := 1
	for i := 1; i <= size; i++ {
		go GetParseUsers("https://www.avito.ru/rostov-na-donu/kvartiry?p="+strconv.Itoa(i), &result, group)
	}
	for endItem <= size {
		select {
		case <-group:
			endItem++
		}
	}
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
}

func GetParseUsers(url string, result *[]FindItem, group chan int) (bool, error) {
	res, err := http.Get(url)
	if err != nil {
		group <- 0
		return false, err
	}
	defer res.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(res.Body)
	doc.Find(".item").Each(func(i int, s *goquery.Selection) {
		id, okId := s.Attr("data-item-id")
		if okId == true {
			idInt, _ := strconv.Atoi(id)
			findItem := findItemInResult(idInt, result)
			if findItem == false && s.Find(".item-description-title span").Text() != "" {
				newItem := FindItem{
					Title:    s.Find(".item-description-title span").Text(),
					Price:    s.Find(".about .price").Text(),
					Location: s.Find(".addres").Text(),
					Id:       idInt,
				}
				*result = append(*result, newItem)
			}
		}
	})
	group <- 1
	return true, nil
}

func findItemInResult(item int, result *[]FindItem) bool {
	for i := 0; i < len(*result); i++ {
		if (*result)[i].Id == item {
			return true
		}
	}
	return false
}

func addItem(i int, result *[]FindItem, group chan int) {
	GetParseUsers("https://www.avito.ru/rostov-na-donu/kvartiry?p="+strconv.Itoa(i), result, group)
}
