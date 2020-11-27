package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Product struct {
	ProductID      int    `json:"productId"`
	Manufacturer   string `json:"manufacturer"`
	Sku            string `json:"sku"`
	Upc            string `json:"upc"`
	PricePerUnit   string `json:"pricePerUnit"`
	QuantityOnHand int    `json:"quantityOnHand"`
	ProductName    string `json:"productName"`
}

var productList []Product

//you can rename the fields when you're encoding this object to json using 'tags'
//omitempty - remove the field from json, if it's zero or nil
type foo struct {
	Message string `json:"message,omitempty"` //no space, no warning!
	Age     int    `json:"age,omitempty"`
	Name    string `json:"firstName,omitempty"`
	Surname string `json:"lastName,omitempty"`
}

// type fooHandler struct {
// 	Message string
// }

// func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte(f.Message))
// }

// func barHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("HeyaHeya"))
// }

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productsJson, err := json.Marshal(productList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productsJson)
	case http.MethodPost:
		var newProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &newProduct)
		if err != nil || newProduct.ProductID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		newProduct.ProductID = getNextID()
		productList = append(productList, newProduct) // append returns a new slice, that's why we need to re-assign
		w.WriteHeader(http.StatusCreated)
		return
	}
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "products/") // get the text that comes directly after products
	productID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	product, listItemIndex := findProductByID(productID)
	if product == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//now that we have the item and its index, its time to determine the method (get or put(for update))
	switch r.Method {
	case http.MethodGet:
		//return a single product
		productJSON, err := json.Marshal(product)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(productJSON)

	case http.MethodPut:
	}
}

func init() {
	productsJSON := `[
		{
			"productId": 1,
			"manufacturer": "Johns-Jenkins",
			"sku": "p5z343vdS",
			"upc": "939581000000",
			"pricePerUnit": "497.45",
			"quantityOnHand": 9703,
			"productName": "sticky note"
		  },
		  {
			"productId": 2,
			"manufacturer": "Hessel, Schimmel and Feeney",
			"sku": "i7v300kmx",
			"upc": "740979000000",
			"pricePerUnit": "282.29",
			"quantityOnHand": 9217,
			"productName": "leg warmers"
		  },
		  {
			"productId": 3,
			"manufacturer": "Swaniawski, Bartoletti and Bruen",
			"sku": "q0L657ys7",
			"upc": "111730000000",
			"pricePerUnit": "436.26",
			"quantityOnHand": 5905,
			"productName": "lamp shade"
		  },
		  {
			"productId": 4,
			"manufacturer": "Runolfsdottir, Littel and Dicki",
			"sku": "x78426lq1",
			"upc": "93986215015",
			"pricePerUnit": "537.90",
			"quantityOnHand": 2642,
			"productName": "flowers"
		  },
		  {
			"productId": 5,
			"manufacturer": "Kuhn, Cronin and Spencer",
			"sku": "r4X793mdR",
			"upc": "260149000000",
			"pricePerUnit": "112.10",
			"quantityOnHand": 6144,
			"productName": "clamp"
		  }
	]`

	err := json.Unmarshal([]byte(productsJSON), &productList)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	http.HandleFunc("/products", productsHandler)
	http.HandleFunc("/products/", productHandler)
	http.ListenAndServe(":5000", nil)
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

//temp solution for the productID before the DB is implemented:
func getNextID() int {
	highestID := -1
	for _, product := range productList {
		if highestID < product.ProductID {
			highestID = product.ProductID
		}
	}
	return highestID + 1
}

func findProductByID(productID int) (*Product, int) {
	for i, product := range productList {
		if productID == product.ProductID {
			return &product, i
		}
	}
	return nil, 0
}
