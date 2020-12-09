/*Data access code, layer*/
package product

import (
	"Webservice/database"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

//it'll allow us to access product w/o interating through the entire slice
// var productMap = struct {
//embedded struct:
// sync.RWMutex //because webservices are multi-threaded, maps in go are inherently not thread-safe,
//which means we gotta wrap it in mutex to avoid 2 threads from reading and writing in it at the same time
// 	m map[int]Product
// }{m: make(map[int]Product)}

//load the data from json file:
// func init() {
// 	fmt.Println("loading products...")
// 	prodMap, err := loadProductMap()

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	productMap.m = prodMap
// 	fmt.Printf("%d products loaded...\n", len(productMap.m))
// }

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

func getProduct(productID int) (*Product, error) {
	// productMap.RLock()         //prevents other threads from writing
	// defer productMap.RUnlock() //releases the lock
	// //if another thread had a writelock then our code would wait until it's released

	// //if it returns a value
	// if product, ok := productMap.m[productID]; ok {
	// 	return &product
	// }

	// return nil
	var p Product
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT productId, manufacturer, sku, upc, pricePerUnit,
	quantityOnHand, productName FROM products WHERE productId = ?`, productID)
	err := row.Scan(&p.ProductID, &p.Manufacturer, &p.Sku, &p.Upc, &p.PricePerUnit, &p.QuantityOnHand, &p.ProductName)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &p, nil
}

func removeProduct(productID int) error {
	// productMap.Lock()         //locks for writing
	// defer productMap.Unlock() //allow writing
	// delete(productMap.m, productID)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM products WHERE productId=?`, productID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//return a full list of products as a slice
func getProductList() ([]Product, error) {
	// productMap.RLock()
	// products := make([]Product, 0, len(productMap.m)) //size 0, capacity length of our map
	// for _, value := range productMap.m {
	// 	products = append(products, value) //append from products value
	// }
	// productMap.RUnlock()
	// return products
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT productId, manufacturer, sku, upc, pricePerUnit,
	quantityOnHand, productName FROM products`)
	if err != nil {
		log.Fatal(err)
	}
	defer results.Close() //make the db connection available for other queries
	products := make([]Product, 0)
	for results.Next() {
		var p Product
		err = results.Scan(&p.ProductID, &p.Manufacturer, &p.Sku, &p.Upc, &p.PricePerUnit, &p.QuantityOnHand, &p.ProductName)
		products = append(products, p)
	}
	//in case any sort of error occurs in the for loop, we should log it
	err = results.Err()
	if err != nil {
		log.Fatal(err)
	}

	return products, nil
}

func GetTopTenProducts() ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
	productId, 
	manufacturer, 
	sku, 
	upc, 
	pricePerUnit, 
	quantityOnHand, 
	productName 
	FROM products ORDER BY quantityOnHand DESC LIMIT 10
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	products := make([]Product, 0)
	for results.Next() {
		var product Product
		results.Scan(&product.ProductID,
			&product.Manufacturer,
			&product.Sku,
			&product.Upc,
			&product.PricePerUnit,
			&product.QuantityOnHand,
			&product.ProductName)

		products = append(products, product)
	}
	return products, nil
}

//sort in ascending order
// func getProductIds() []int {
// 	productMap.RLock()
// 	productIds := []int{}
// 	for key := range productMap.m {
// 		productIds = append(productIds, key)
// 	}
// 	productMap.RUnlock()
// 	sort.Ints(productIds)
// 	return productIds
// }

//get the next highest ID value
// func getNextProductID() int {
// 	productIds := getProductIds()
// 	return productIds[len(productIds)-1] + 1
// }

// func addOrUpdateProduct(product Product) (int, error) {
// 	// if the product id is set, update, otherwise add
// 	addOrUpdateID := -1
// 	if product.ProductID > 0 {
// 		oldProduct, _ := getProduct(product.ProductID)
// 		// if it exists, replace it, otherwise return error
// 		if oldProduct == nil {
// 			return 0, fmt.Errorf("product id [%d] doesn't exist", product.ProductID)
// 		}
// 		addOrUpdateID = product.ProductID
// 	} else {
// 		addOrUpdateID = getNextProductID()
// 		product.ProductID = addOrUpdateID
// 	}
// 	productMap.Lock()
// 	productMap.m[addOrUpdateID] = product
// 	productMap.Unlock()
// 	return addOrUpdateID, nil
// }

func updateProduct(product Product) error {
	if product.ProductID == 0 || &product == nil {
		return errors.New("product has invalid ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `UPDATE products SET manufacturer=?, 
	sku=?, 
	upc=?, 
	pricePerUnit=CAST(? AS DECIMAL(13,2)), 
	quantityOnHand=?, 
	productName=? WHERE productId=?`, product.Manufacturer, product.Sku, product.Upc, product.PricePerUnit, product.QuantityOnHand, product.ProductName, product.ProductID)

	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func insertProduct(product Product) (int, error) {
	if &product == nil {
		return 0, errors.New("product is invalid")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := database.DbConn.ExecContext(ctx, `INSERT INTO products (manufacturer, sku, upc, pricePerUnit, quantityOnHand, productName)
	VALUES=(?,?,?,?,?, ?)`, product.Manufacturer, product.Sku, product.Upc, product.PricePerUnit, product.QuantityOnHand, product.ProductName)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	return int(insertId), nil
}

func searchForProductData(productFilter ProductReportFilter) ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var queryArgs = make([]interface{}, 0)
	var queryBuilder strings.Builder
	queryBuilder.WriteString(`SELECT 
		productId, 
		LOWER(manufacturer), 
		LOWER(sku), 
		upc, 
		pricePerUnit, 
		quantityOnHand, 
		LOWER(productName) 
		FROM products WHERE `)
	if productFilter.NameFilter != "" {
		queryBuilder.WriteString(`productName LIKE ? `)
		queryArgs = append(queryArgs, "%"+strings.ToLower(productFilter.NameFilter)+"%")
	}
	if productFilter.ManufacturerFilter != "" {
		if len(queryArgs) > 0 {
			queryBuilder.WriteString(" AND ")
		}
		queryBuilder.WriteString(`manufacturer LIKE ? `)
		queryArgs = append(queryArgs, "%"+strings.ToLower(productFilter.ManufacturerFilter)+"%")
	}
	if productFilter.SKUFilter != "" {
		if len(queryArgs) > 0 {
			queryBuilder.WriteString(" AND ")
		}
		queryBuilder.WriteString(`sku LIKE ? `)
		queryArgs = append(queryArgs, "%"+strings.ToLower(productFilter.SKUFilter)+"%")
	}

	results, err := database.DbConn.QueryContext(ctx, queryBuilder.String(), queryArgs...)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	products := make([]Product, 0)
	for results.Next() {
		var product Product
		results.Scan(&product.ProductID,
			&product.Manufacturer,
			&product.Sku,
			&product.Upc,
			&product.PricePerUnit,
			&product.QuantityOnHand,
			&product.ProductName)

		products = append(products, product)
	}
	return products, nil
}
