package parseUsers

import (
	"log"
	"net/http"
)

type FindUser struct {
	Name, UrlProf, Avatar string
}

func GetParseUsers(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
}

func addUsers(i int) {
	go GetParseUsers("https://www.twitch.tv/qadRat")
}

func Start(size int) {
	for i := 1; i < size; i++ {
		go addUsers(i)
	}
}