package receipt

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

var ReceiptDirectory string = filepath.Join("uploads")

type Receipt struct {
	ReceiptName string    `json:"name"`
	UploadDate  time.Time `json:"uploadDate"`
}

func GetReceipts() ([]Receipt, error) {
	receipts := []Receipt{}
	files, err := os.ReadDir(ReceiptDirectory)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		var receipt Receipt
		fileModTime, err := os.Stat(filepath.Join(ReceiptDirectory, file.Name()))
		if err != nil {
			log.Fatal(err)
		}
		receipt.ReceiptName = file.Name()
		receipt.UploadDate = fileModTime.ModTime()

		receipts = append(receipts, receipt)
	}

	return receipts, nil
}
