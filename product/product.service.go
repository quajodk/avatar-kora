package product

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	product, err := getProduct(productID)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
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
		if newProduct.ProductID != productID {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		err = updateProduct(newProduct)
		if err != nil {
			log.Fatal(err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
		return
	case http.MethodDelete:
		err := removeProduct(productID)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	case http.MethodOptions:
		return
	default:
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func productsHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		productList, err := getProductList()
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		productsJson, err := json.Marshal(productList)
		if err != nil {
			log.Fatal(err)
		}
		res.Header().Set("Content-Type", "application/json")
		_, err = res.Write(productsJson)
		if err != nil {
			log.Fatal(err)
		}
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
		id, err := addProduct(newProduct)
		if err != nil {
			log.Fatal(err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(id)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusCreated)
		res.Write(j)
		return
	case http.MethodOptions:
		return
	}
}
