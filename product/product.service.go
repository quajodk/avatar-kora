package product

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"services/cors"
	"strconv"
	"strings"
)

const productsBasePath = "products"

func SetupRoutes(apiBasePath string) {
	handleProducts := http.HandlerFunc(productsHandler)
	handleProduct := http.HandlerFunc(productHandler)

	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, productsBasePath), cors.Middleware(handleProducts))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, productsBasePath), cors.Middleware(handleProduct))
}

func productHandler(res http.ResponseWriter, req *http.Request) {
	urlPathSeg := strings.Split(req.URL.Path, "products/")
	productID, err := strconv.Atoi(urlPathSeg[len(urlPathSeg)-1])
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	product := getProduct(productID)
	if product == nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	switch req.Method {
	case http.MethodGet:
		// return a single product
		productJSON, err := json.Marshal(product)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.Write(productJSON)
	case http.MethodPut:
		// update product in the list
		var updateProduct Product
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &updateProduct)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if updateProduct.ProductID != productID {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		addOrUpdateProduct(updateProduct)
		res.WriteHeader(http.StatusOK)
		return
	case http.MethodDelete:
		removeProduct(productID)
	case http.MethodOptions:
		return
	default:
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func productsHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		productList := getProductList()
		productsJson, err := json.Marshal(productList)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.Write(productsJson)
	case http.MethodPost:
		// add a new product to the list
		var newProduct Product
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &newProduct)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if newProduct.ProductID != 0 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = addOrUpdateProduct(newProduct)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
		}

		res.WriteHeader(http.StatusCreated)
		return
	case http.MethodOptions:
		return
	}
}
