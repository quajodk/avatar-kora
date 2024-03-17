package receipt

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"services/cors"
	"strconv"
	"strings"
)

const receiptPath = "receipts"

func handleReceipts(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		receiptList, err := GetReceipts()
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		receiptJson, err := json.Marshal(receiptList)
		if err != nil {
			log.Fatal(err)
		}
		_, err = res.Write(receiptJson)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodPost:
		req.ParseMultipartForm(5 << 20) // 5mb
		file, handler, err := req.FormFile("receipt")
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		defer file.Close()
		f, err := os.OpenFile(filepath.Join(ReceiptDirectory, handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		res.WriteHeader(http.StatusCreated)
	case http.MethodOptions:
		return
	default:
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func handleDownload(res http.ResponseWriter, req *http.Request) {
	urlPathSeg := strings.Split(req.URL.Path, fmt.Sprintf("%s/", receiptPath))
	if len(urlPathSeg[1:]) > 1 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	fileName := urlPathSeg[1:][0]
	file, err := os.Open(filepath.Join(ReceiptDirectory, fileName))
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()
	fileHeader := make([]byte, 512)
	file.Read(fileHeader)
	fileContentType := http.DetectContentType(fileHeader)
	stat, err := file.Stat()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	fileSize := strconv.FormatInt(stat.Size(), 10)
	res.Header().Add("Content-Disposition", "attachment; filename="+fileName)
	res.Header().Add("Content-Type", fileContentType)
	res.Header().Add("Content-Length", fileSize)
	file.Seek(0, 0)
	io.Copy(res, file)
}

func SetupRoutes(apiBasePath string) {
	receiptHandler := http.HandlerFunc(handleReceipts)
	downloadHandler := http.HandlerFunc(handleDownload)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, receiptPath), cors.Middleware(receiptHandler))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, receiptPath), cors.Middleware(downloadHandler))

}
