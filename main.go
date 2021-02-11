package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
)

type Image struct {
	ID        string
	Imagepath string
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
		panic(err)
	}
	fmt.Println("-----------------------------------------------")
}
func form(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("form.html"))

	tmpl.Execute(w, nil)

}
func upload(w http.ResponseWriter, r *http.Request) {

	file, handler, err := r.FormFile("myFile")

	if err != nil {
		panic(err)
	}
	defer file.Close()
	imagePath := filepath.Join("store/", handler.Filename)
	dest, err := os.Create(imagePath)
	if err != nil {
		panic(err)
	}
	defer dest.Close()
	io.Copy(dest, file)
	sqlStatement := "INSERT INTO images (imagepath) VALUES (?)"
	_, err = db.Exec(sqlStatement, imagePath)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/", 301)

}
func display(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("display.html")
	if err != nil {
		fmt.Println(err)
	}
	var image []Image

	query := "SELECT *FROM images"
	rows, err := db.Query(query)
	for rows.Next() {
		var img Image
		err = rows.Scan(&img.ID, &img.Imagepath)
		if err != nil {
			panic(err)
		}
		image = append(image, img)
	}
	tmpl.Execute(w, image)
}
func main() {
	http.HandleFunc("/", form)
	http.HandleFunc("/display", display)
	http.HandleFunc("/upload", upload)
	http.ListenAndServe(":8080", nil)
}
