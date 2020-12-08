/*By this argument, PUT is for creating when you know the URL of the thing you will create.
POST can be used to create when you know the URL of the "factory" or manager for the category of things you want to create.
source: https://stackoverflow.com/questions/630453/put-vs-post-in-rest#:~:text=You%20can%20PUT%20a%20resource,the%20thing%20you%20will%20create.
*/
package product

import (
	"Webservice/cors"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/websocket"
)

const productsPath = "products"

// convert the handler functions into handlers and register them for specific routes
func SetupRoutes(apiBasePath string) {
	//teacher:
	// productsHandler := http.HandlerFunc(handleProducts)
	// productHandler := http.HandlerFunc(handleProduct)
	// http.Handle(fmt.Sprintf("%s/%s", apiBasePath, productsPath), cors.Middleware(productsHandler))
	// http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, productsPath), cors.Middleware(productHandler))
	//me:
	http.HandleFunc(fmt.Sprintf("%s/%s", apiBasePath, productsPath), cors.MiddlewareFunc(productsHandler))
	http.HandleFunc(fmt.Sprintf("%s/%s/", apiBasePath, productsPath), cors.MiddlewareFunc(productHandler))
	http.Handle("/websocket", websocket.Handler(productSocket))
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productList, err := getProductList()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
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

		_, err = insertProduct(newProduct)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// productList = append(productList, newProduct) // append returns a new slice, that's why we need to re-assign
		w.WriteHeader(http.StatusCreated)
		return
	/*part of CORS workflow:
	browser sends a special type of request - pre-fly request
	that uses http.OptionsMethod. Webservice then returns CORS specific headers
	so that the browser knows if it should allow the traffic to be sent to that server*/
	case http.MethodOptions:
		//we just return here because our middleware handles return of the headers for us
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "products/") // get the text that comes directly after products
	productID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// product, listItemIndex := findProductByID(productID)
	product, err := getProduct(productID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
		var updatedProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &updatedProduct)
		if err != nil || updatedProduct.ProductID != productID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// product = &updatedProduct
		// productList[listItemIndex] = *product
		updateProduct(updatedProduct)
		w.WriteHeader(http.StatusOK)
	case http.MethodOptions:
		return
	case http.MethodDelete:
		removeProduct(productID)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	return
}
