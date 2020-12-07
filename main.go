package main

import (
	"Webservice/database"
	"Webservice/product"
	"Webservice/receipt"
	"log"
	"net/http"
)

//you can rename the fields when you're encoding this object to json using 'tags'
//omitempty - remove the field from json, if it's zero or nil
// type foo struct {
// 	Message string `json:"message,omitempty"` //no space, no warning!
// 	Age     int    `json:"age,omitempty"`
// 	Name    string `json:"firstName,omitempty"`
// 	Surname string `json:"lastName,omitempty"`
// }

// type fooHandler struct {
// 	Message string
// }

// func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte(f.Message))
// }

// func barHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("HeyaHeya"))
// }

const basePath = "/api"

func main() {
	database.SetupDatabase()
	product.SetupRoutes(basePath)
	receipt.SetupRoutes(basePath)

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal(err)
	}

	// http.Handle("/foo", &fooHandler{Message: "hello"})
	// http.HandleFunc("/bar", barHandler)
	// http.ListenAndServe(":5000", nil)

	//ENCODING JSON
	// data, _ := json.Marshal(&foo{"4Sore", 56, "Abe", "Lincoln"}) //encodes the object provided into JSON, all object's fields need to be exported (aka capitalized) to show up on jSON
	// fmt.Println(string(data))

	//DECODING JSON
	// f := foo{}
	// err := json.Unmarshal(data, &f)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(f.Message)

	//OR

	// foo := func(w http.ResponseWriter, _ *http.Request) {
	// 	w.Write([]byte("Hello there"))
	// }
	// http.HandleFunc("/foo", foo)
}
