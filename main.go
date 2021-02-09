package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "UserDB"
)

type Data struct {
	ID        string
	FirstName string
	LastName  string
	DOB       string
	Email     string
	Phone     string
}

var db *sql.DB

func init() {
	pInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", pInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic((err))
	}

	fmt.Println("---------------------------------------------------")

}
func insert(w http.ResponseWriter, r *http.Request) {
	sqlStatement := "INSERT INTO userdata(fname, lname, dob, email, phone) VALUES ($1, $2, $3, $4, $5)"
	inputData := Data{
		FirstName: r.FormValue("fname"),
		LastName:  r.FormValue("lname"),
		DOB:       r.FormValue("dob"),
		Email:     r.FormValue("email"),
		Phone:     r.FormValue("phone"),
	}

	_, err := db.Exec(sqlStatement, inputData.FirstName, inputData.LastName, inputData.DOB, inputData.Email, inputData.Phone)
	fmt.Println("done")
	http.Redirect(w, r, "/", 301)
	if err != nil {

		fmt.Println(err)
	}
}
func form(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("form.html"))
	tmpl.Execute(w, nil)

}
func display(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("display.html"))

	rows, err := db.Query("SELECT * FROM userdata")
	if err != nil {
		fmt.Println(err)
		fmt.Println("erorrrrr")
	}
	var getAll []Data

	for rows.Next() {

		var get Data
		err = rows.Scan(&get.FirstName, &get.LastName, &get.DOB, &get.Email, &get.Phone, &get.ID)
		if err != nil {
			fmt.Println(err)
		}
		getAll = append(getAll, get)

	}

	tmpl.Execute(w, getAll)

}
func edit(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("edit.html"))

	fmt.Println("here")
	id := r.URL.Query().Get("id")

	var get Data
	sqlStatement := "SELECT * FROM userdata where id =" + id

	err := db.QueryRow(sqlStatement).Scan(&get.FirstName, &get.LastName, &get.DOB, &get.Email, &get.Phone, &get.ID)

	if err != nil {
		fmt.Println(err)
	}
	if r.Method != http.MethodPost {
		tmpl.Execute(w, get)

		return
	}
	query := "UPDATE userdata SET fname = $1, lname = $2, dob = $3, email = $4, phone = $5 WHERE id = $6"
	inputData := Data{
		FirstName: r.FormValue("fname"),
		LastName:  r.FormValue("lname"),
		DOB:       r.FormValue("dob"),
		Email:     r.FormValue("email"),
		Phone:     r.FormValue("phone"),
	}
	_, err = db.Exec(query, inputData.FirstName, inputData.LastName, inputData.DOB, inputData.Email, inputData.Phone, id)

	http.Redirect(w, r, "/display", 301)
}

func delete(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	sqlStatement := "DELETE FROM userdata WHERE id = $1"
	_, err := db.Exec(sqlStatement, id)
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/display", 301)

}
func main() {

	http.HandleFunc("/", form)
	http.HandleFunc("/insert", insert)
	http.HandleFunc("/display", display)
	http.HandleFunc("/edit", edit)
	http.HandleFunc("/delete", delete)
	http.ListenAndServe(":8080", nil)
}
