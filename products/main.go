package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type products struct {
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
		ID             int    `json:"id"`
		Name           string `json:"name"`
		Description    string `json:"description"`
		Image          string `json:"image"`
		Price          string `json:"price"`
		DiscountAmount string `json:"discount_amount"`
		Status         bool   `json:"status"`
		Categories     []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"categories"`
	} `json:"data"`
}

func main() {
	url := "https://gorest.co.in/public-api/products"

	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	info, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	var all products
	err = json.Unmarshal([]byte(info), &all)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(all)

	file, err := json.MarshalIndent(all, "", "")
	if err != nil {
		fmt.Println(err)
	}

	err = ioutil.WriteFile("products.json", file, 0644)
	if err != nil {
		fmt.Println(err)
	}

}
