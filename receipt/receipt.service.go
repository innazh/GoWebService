package receipt

import (
	"Webservice/cors"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const receiptPath = "receipts"

func handleReceipts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		receiptList, err := GetReceipts()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		receiptListJson, err := json.Marshal(receiptList)
		if err != nil {
			log.Fatal(err)
		}
		_, err = w.Write(receiptListJson)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodPost:
		r.ParseMultipartForm(5 << 20) //5MB ???
		recFile, fheader, err := r.FormFile("receipt")
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer recFile.Close()
		// prob need to open before the copy
		f, err := os.OpenFile(filepath.Join(ReceiptDirectory, fheader.Filename), os.O_WRONLY|os.O_CREATE, 0666) //todo: read abt permissions in-depth on unix again
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer f.Close()
		//copy the received through the post method file to the file that we created in out local dir
		io.Copy(f, recFile)
		w.WriteHeader(http.StatusCreated)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func SetupRoutes(apiBasePath string) {
	http.HandleFunc(fmt.Sprintf("%s/%s", apiBasePath, receiptPath), cors.MiddlewareFunc(handleReceipts))
}
