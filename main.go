package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

const (
	driver   = "mysql"
	user     = "root"
	password = "Charmi@123"
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
	var err error
	db, err = sql.Open("mysql", "superuser1:Super_101@(127.0.0.1:3306)/UserDB")
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

	sqlStatement := "INSERT INTO userdata(fname, lname, dob, email, phone) VALUES (?,?,?,?,?)"
	inputData := Data{
		FirstName: r.FormValue("fname"),
		LastName:  r.FormValue("lname"),
		DOB:       r.FormValue("dob"),
		Email:     r.FormValue("email"),
		Phone:     r.FormValue("phone"),
	}
	stmt, err := db.Prepare(sqlStatement)
	if err != nil {
		fmt.Println(err)
	}
	_, err = stmt.Exec(inputData.FirstName, inputData.LastName, inputData.DOB, inputData.Email, inputData.Phone)

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

	}
	var getAll []Data
	for rows.Next() {

		var get Data
		err = rows.Scan(&get.ID, &get.FirstName, &get.LastName, &get.DOB, &get.Email, &get.Phone)
		if err != nil {
			fmt.Println(err)
		}
		getAll = append(getAll, get)

	}

	tmpl.Execute(w, getAll)

}
func edit(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("edit.html"))
	id := r.URL.Query().Get("id")

	var get Data
	sqlStatement := "SELECT * FROM userdata where id =" + id

	rows := db.QueryRow(sqlStatement)
	err := rows.Scan(&get.ID, &get.FirstName, &get.LastName, &get.DOB, &get.Email, &get.Phone)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(get)
	if r.Method != http.MethodPost {
		tmpl.Execute(w, get)

		return
	}
	query := "UPDATE userdata SET fname = ?, lname = ?, dob = ?, email = ?, phone = ? WHERE id = ?"
	inputData := Data{
		FirstName: r.FormValue("fname"),
		LastName:  r.FormValue("lname"),
		DOB:       r.FormValue("dob"),
		Email:     r.FormValue("email"),
		Phone:     r.FormValue("phone"),
	}

	stmt, err := db.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	_, err = stmt.Exec(inputData.FirstName, inputData.LastName, inputData.DOB, inputData.Email, inputData.Phone, id)
	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, "/display", 301)
}

func delete(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	sqlStatement := "DELETE FROM userdata WHERE id = ?"
	stmt, err := db.Prepare(sqlStatement)
	if err != nil {
		fmt.Println(err)
	}
	_, err = stmt.Exec(id)
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
