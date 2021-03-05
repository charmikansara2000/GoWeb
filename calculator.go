package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Expression struct {
	Exp    string `json:"exp"`
	Result string `json:"result"`
}

func FetchAndPerform(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var expression Expression
	json.NewDecoder(r.Body).Decode(&expression)

	fmt.Println(expression)
	var ans int
	s := expression.Exp
	for i := range s {
		char := string(s[i])
		if char == "+" || char == "-" || char == "/" || char == "*" {
			res := strings.Split(s, char)
			fmt.Println(res)
			a, _ := strconv.Atoi(res[0])
			b, _ := strconv.Atoi(res[1])

			switch char {
			case "+":
				ans = a + b
			case "-":
				ans = a - b
			case "/":
				ans = a / b
			case "*":
				ans = a * b
			}
			//ans = a + b
			fmt.Println(ans)

		}
	}
	expression.Result = fmt.Sprint(ans)
	resp, _ := json.Marshal(expression)
	_ = ioutil.WriteFile("history.json", resp, 0644)
	w.Write(resp)

	return

}

// func Values(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "POST" {

// 		inputValues := values{
// 			op1:     r.FormValue("bttn1"),
// 			op2:     r.FormValue("bttn2"),
// 			operand: r.FormValue("bttnplus"),
// 		}
// 		fmt.Println(inputValues)

// 	}

// }
func main() {
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "Welcome to my website!")
	// })

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/fetch", FetchAndPerform)
	//	http.HandleFunc("/Values", Values)
	http.ListenAndServe(":9090", nil)
}
