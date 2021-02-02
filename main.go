package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type Data struct {
	FirstName string
	LastName  string
	DOB       string
	Email     string
	Phone     string
}

func openForm(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("form.html")
	tmpl.Execute(w, nil)
}
func addData(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		userData, err := ioutil.ReadFile("form.json")
		if err != nil {
			fmt.Println(err)
		}

		var allData []Data
		err = json.Unmarshal([]byte(userData), &allData)
		if err != nil {
			fmt.Println("Error in unmarshalling data from file")
			fmt.Println(err)
		}

		//var dt Data
		inputData := Data{
			FirstName: r.FormValue("fname"),
			LastName:  r.FormValue("lname"),
			DOB:       r.FormValue("dob"),
			Email:     r.FormValue("email"),
			Phone:     r.FormValue("phone"),
		}

		allData = append(allData, inputData)
		file, _ := json.MarshalIndent(allData, "", "")
		_ = ioutil.WriteFile("form.json", file, 0644)

		tmpl, _ := template.ParseFiles("thankyou.html")
		tmpl.Execute(w, nil)

	}
}
func main() {
	http.HandleFunc("/", openForm)
	http.HandleFunc("/form", addData)
	http.ListenAndServe(":8080", nil)
}
