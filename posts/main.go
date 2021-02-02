package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type posts struct {
	Code int `json:"code"`
	Meta struct {
		Pagination struct {
			Total int `json:"total"`
			Pages int `json:"pages"`
			Page  int `json:"page"`
			Limit int `json:"limit"`
		} `json:"pagination"`
	} `json:"meta"`
	Data []struct {
		ID        int       `json:"id"`
		UserID    int       `json:"user_id"`
		Title     string    `json:"title"`
		Body      string    `json:"body"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"data"`
}

func main() {
	url := "https://gorest.co.in/public-api/posts"

	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	info, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	var all posts
	err = json.Unmarshal([]byte(info), &all)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(all)

	file, err := json.MarshalIndent(all, "", "")
	if err != nil {
		fmt.Println(err)
	}

	err = ioutil.WriteFile("posts.json", file, 0644)
	if err != nil {
		fmt.Println(err)
	}

}
