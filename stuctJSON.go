package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type ADDRESS struct {
	Area    string `json:"area"`
	Country string `json:"country"`
}
type User struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Address ADDRESS `json:"address"`
}
type TECHDETS struct {
	Tech string  `json:"tech"`
	Exp  float64 `json:"exp"`
}
type Tech struct {
	Id       int        `json:"id"`
	Techdets []TECHDETS `json:"techdets"`
}
type ContactDETS struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}
type Contact struct {
	Id          int         `json:"id"`
	ContactDets ContactDETS `json:"cotactDets"`
}
type ALL struct {
	//id int
	Name     string
	Address  ADDRESS
	Techdets []TECHDETS
	Email    string
	Phone    string
}

func main() {

	userJson, err := os.Open("userS1.json")
	if err != nil {
		fmt.Println(err)
	}
	userInBytes, _ := ioutil.ReadAll(userJson)
	var users []User
	json.Unmarshal([]byte(userInBytes), &users)

	contactJson, err := os.Open("contactS2.json")
	if err != nil {
		fmt.Println(err)
	}
	contactInBytes, _ := ioutil.ReadAll(contactJson)
	var contacts []Contact
	json.Unmarshal([]byte(contactInBytes), &contacts)

	techJson, err := os.Open("techS3.json")
	if err != nil {
		fmt.Println(err)
	}
	techInBytes, _ := ioutil.ReadAll(techJson)
	var techs []Tech
	json.Unmarshal([]byte(techInBytes), &techs)

	MAP := make(map[string]string)
	MAP["IND"] = "+91"
	MAP["UK"] = "+41"

	answer := make(map[string]ALL, 2)

	for _, user := range users {
		var temp ALL

		id := user.Id
		temp.Name = user.Name
		temp.Address = user.Address
		country := user.Address.Country
		for _, tech := range techs {
			if tech.Id == id {
				temp.Techdets = temp.Techdets
			}
		}
		for _, contact := range contacts {
			if contact.Id == id {
				temp.Email = contact.ContactDets.Email
				temp.Phone = contact.ContactDets.Phone
			}
		}
		for j := range MAP {
			if j == country {
				temp.Phone = MAP[j] + temp.Phone
			}
		}
		answer[strconv.Itoa(id)] = temp
	}

	fmt.Println(answer)
}

/*s1 := make([]User, 2)
s1[0] = User{id: 1, name: "radha", address: ADDRESS{area: "bopal", country: "IND"}}
s1[1] = User{id: 2, name: "aniket", address: ADDRESS{area: "maninagar", country: "UK"}}

s2 := make([]Contact, 2)
s2[0] = Contact{id: 1, techdets: []TECHDETS{{tech: "react", exp: 1}, {tech: "go", exp: 1.5}}}
s2[1] = Contact{id: 2, techdets: []TECHDETS{{tech: "react", exp: 0.9}, {tech: "go", exp: 1.5}}}

s3 := make([]Contact, 2)
s3[0] = Contact{id: 1, contactDets: ContactDETS{email: "abc@gmail.com", phone: "123455778"}}
s3[1] = Contact{id: 2, contactDets: ContactDETS{email: "xyz@gmail.com", phone: "123455778"}}*/
