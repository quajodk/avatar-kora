package product

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"services/database"
	"sync"
	"time"
)

// used to store product in memory
var productMap = struct {
	sync.RWMutex
	m map[int]Product
}{m: make(map[int]Product)}

func init() {
	fmt.Println("loading products ...")
	prodMap, err := loadProductMap()
	productMap.m = prodMap
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d products loaded... \n", len(productMap.m))
}

func loadProductMap() (map[int]Product, error) {
	fileName := "products.json"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file [%s] does not exist", fileName)
	}
	file, _ := os.ReadFile(fileName)
	productList := make([]Product, 0)
	err = json.Unmarshal([]byte(file), &productList)
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
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	result := database.DB.QueryRowContext(ctx, "SELECT productid, manufacturer, sku, upc, priceperunit, quantityonhand, productname FROM products WHERE productid = $1", productID)
	var product Product
	err := result.Scan(&product.ProductID, &product.Manufacturer, &product.Sku, &product.Upc, &product.PricePerUnit, &product.QuantityOnHand, &product.ProductName)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &product, nil

}

func removeProduct(productID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := database.DB.QueryContext(ctx, "DELETE FROM products WHERE productid = $1", productID)
	if err != nil {
		return err
	}
	return nil
}

func getProductList() ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	result, err := database.DB.QueryContext(ctx, `SELECT productid, manufacturer, sku, upc, priceperunit, quantityonhand, productname FROM products ORDER BY productid ASC`)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	products := make([]Product, 0)
	for result.Next() {
		var product Product
		result.Scan(&product.ProductID, &product.Manufacturer, &product.Sku, &product.Upc, &product.PricePerUnit, &product.QuantityOnHand, &product.ProductName)
		products = append(products, product)
	}

	return products, nil
}

func updateProduct(product Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := database.DB.ExecContext(ctx, "UPDATE products SET manufacturer = $1, sku = $2, upc = $3, priceperunit = $4, quantityonhand = $5, productname = $6 WHERE productid = $7", product.Manufacturer, product.Sku, product.Upc, product.PricePerUnit, product.QuantityOnHand, product.ProductName, product.ProductID)
	if err != nil {
		return err
	}
	return nil
}

func addProduct(product Product) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var productID int64
	err := database.DB.QueryRowContext(ctx, "INSERT INTO products(manufacturer, sku, upc, priceperunit, quantityonhand, productname) VALUES($1, $2, $3, $4, $5, $6) RETURNING productid", product.Manufacturer, product.Sku, product.Upc, product.PricePerUnit, product.QuantityOnHand, product.ProductName).Scan(&productID)
	if err != nil {
		return 0, err
	}
	id := productID
	return int(id), nil
}
