package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Credentials struct {
	gorm.Model
	Id       int
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
type Token struct {
	TokenString string `json:"tokenstring"`
}
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

var jwtKey = []byte("my_key")
var db *gorm.DB

func migration() {
	var err error
	db, err = gorm.Open(postgres.Open("postgres://postgres:postgres@localhost/UserDB?sslmode=disable"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	db.AutoMigrate(&Credentials{})
	fmt.Println("----------------------------")

}
func SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user Credentials
	json.NewDecoder(r.Body).Decode(&user)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	user.Password = string(hashedPassword)
	if err != nil {
		fmt.Println(err)
	}
	db.Create(&user)
	json.NewEncoder(w).Encode(user)
}
func SignIn(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var cred Credentials

	err := json.NewDecoder(r.Body).Decode(&cred)

	if err != nil {
		fmt.Println("AAAAAAA")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var storedCreds Credentials
	db.Where("Username = ?", cred.Username).First(&storedCreds)

	expirationTime := time.Now().Add(50 * time.Minute)

	claims := &Claims{
		Username: cred.Username,
		Role:     cred.Role,
		StandardClaims: jwt.StandardClaims{

			ExpiresAt: expirationTime.Unix(),
		},
	}
	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(cred.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		fmt.Println("errorrrrrrrr")
		w.WriteHeader(http.StatusUnauthorized)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//var tokens Token
	data := Token{
		TokenString: tokenString,
	}
	bytedata, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	w.Write(bytedata)

}

func Welcome(w http.ResponseWriter, r *http.Request) {
	var tknStr Token
	err := json.NewDecoder(r.Body).Decode(&tknStr)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ts := tknStr.TokenString

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(ts, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Println(claims.Role)
	fmt.Println(claims.Username)
	fmt.Println(claims)
	if claims.Role == "admin" {
		http.Redirect(w, r, "/admin", 301)
		return
	} else if claims.Role == "user" {
		http.Redirect(w, r, "/user", 301)
		return
	} else {
		w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
	}
}
func User(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome User"))
}
func Admin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome Admin"))
}

func main() {

	migration()
	http.HandleFunc("/signup", SignUp)
	http.HandleFunc("/signin", SignIn)
	http.HandleFunc("/welcome", Welcome)
	http.HandleFunc("/user", User)
	http.HandleFunc("/admin", Admin)
	http.ListenAndServe(":9000", nil)
}
