package main

import (
	"net/http"
	"services/database"
	"services/product"

	_ "github.com/lib/pq"
)

const apiBasePath = "/api"

func main() {
	database.IntDB()
	product.SetupRoutes(apiBasePath)
	http.ListenAndServe(":5000", nil)
}
