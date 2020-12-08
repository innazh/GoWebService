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
	"strconv"
	"strings"
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
		r.ParseMultipartForm(5 << 20)                  //5MB //not neccessary
		recFile, fheader, err := r.FormFile("receipt") //sending the request via body->form-data, key=receipt, value=filename of type file (has 2 types: file or text)
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

func handleDownload(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", receiptPath))
	if len(urlPathSegments[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fileName := urlPathSegments[1:][0]
	file, err := os.Open(filepath.Join(ReceiptDirectory, fileName))
	defer file.Close()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fHeader := make([]byte, 512)
	file.Read(fHeader)
	fContentType := http.DetectContentType(fHeader)

	stat, err := file.Stat()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fSize := strconv.FormatInt(stat.Size(), 10)
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", fContentType)
	w.Header().Set("Content-Length", fSize)
	file.Seek(0, 0)
	io.Copy(w, file)

}

func SetupRoutes(apiBasePath string) {
	http.HandleFunc(fmt.Sprintf("%s/%s", apiBasePath, receiptPath), cors.MiddlewareFunc(handleReceipts))
	http.HandleFunc(fmt.Sprintf("%s/%s/", apiBasePath, receiptPath), cors.MiddlewareFunc(handleDownload))
}
