package main

import (
	"net/http"
	"services/database"
	"services/product"
	"services/receipt"

	_ "github.com/lib/pq"
)

const apiBasePath = "/api"

func main() {
	database.IntDB()
	product.SetupRoutes(apiBasePath)
	receipt.SetupRoutes(apiBasePath)
	http.ListenAndServe(":5000", nil)
}
