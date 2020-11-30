/*Data access code, layer*/
package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"
)

//it'll allow us to access product w/o interating through the entire slice
var productMap = struct {
	//embedded struct:
	sync.RWMutex //because webservices are multi-threaded, maps in go are inherently not thread-safe,
	//which means we gotta wrap it in mutex to avoid 2 threads from reading and writing in it at the same time
	m map[int]Product
}{m: make(map[int]Product)}

//load the data from json file:
func init() {
	fmt.Println("loading products...")
	prodMap, err := loadProductMap()

	if err != nil {
		log.Fatal(err)
	}

	productMap.m = prodMap
	fmt.Printf("%d products loaded...\n", len(productMap.m))
}

func loadProductMap() (map[int]Product, error) {
	filename := "products.json"
	_, err := os.Stat(filename) //checks if the file is there, returns the info about the file (which we don't rly need, therefore _)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file [%s] does not exist", filename)
	}

	file, _ := ioutil.ReadFile(filename) //reads all the data in the file into a byte slice
	productList := make([]Product, 0)
	err = json.Unmarshal(file, &productList)

	if err != nil {
		log.Fatal(err)
	}
	prodMap := make(map[int]Product)
	for i := 0; i < len(productList); i++ {
		prodMap[productList[i].ProductID] = productList[i]
	}
	return prodMap, nil
}

func getProduct(productID int) *Product {
	productMap.RLock()         //prevents other threads from writing
	defer productMap.RUnlock() //releases the lock
	//if another thread had a writelock then our code would wait until it's released

	//if it returns a value
	if product, ok := productMap.m[productID]; ok {
		return &product
	}

	return nil
}

func removeProduct(productID int) {
	productMap.Lock()         //locks for writing
	defer productMap.Unlock() //allow writing
	delete(productMap.m, productID)
}

//return a full list of products as a slice
func getProductList() []Product {
	productMap.RLock()
	products := make([]Product, 0, len(productMap.m)) //size 0, capacity length of our map
	for _, value := range productMap.m {
		products = append(products, value) //append from products value
	}
	productMap.RUnlock()
	return products
}

//sort in ascending order
func getProductIds() []int {
	productMap.RLock()
	productIds := []int{}
	for key := range productMap.m {
		productIds = append(productIds, key)
	}
	productMap.RUnlock()
	sort.Ints(productIds)
	return productIds
}

//get the next highest ID value
func getNextProductID() int {
	productIds := getProductIds()
	return productIds[len(productIds)-1] + 1
}

func addOrUpdateProduct(product Product) (int, error) {
	// if the product id is set, update, otherwise add
	addOrUpdateID := -1
	if product.ProductID > 0 {
		oldProduct := getProduct(product.ProductID)
		// if it exists, replace it, otherwise return error
		if oldProduct == nil {
			return 0, fmt.Errorf("product id [%d] doesn't exist", product.ProductID)
		}
		addOrUpdateID = product.ProductID
	} else {
		addOrUpdateID = getNextProductID()
		product.ProductID = addOrUpdateID
	}
	productMap.Lock()
	productMap.m[addOrUpdateID] = product
	productMap.Unlock()
	return addOrUpdateID, nil
}
