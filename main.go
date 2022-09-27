package main

import (
	"emaildata/modelData"
	"emaildata/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const maxUploadSize = 10 * 1024 * 1024 // 10 mb
const uploadPath = "./demo"
const downloadFileFromPath = "demo/"
const destination = "test/download/"

var con = service.Email{}

func init() {
	con.Server = "mongodb://localhost:27017"
	// con.Server = "mongodb+srv://m001-student:m001-mongodb-basics@sandbox.7zffz3a.mongodb.net/?retryWrites=true&w=majority"
	con.Database = "emailData"
	con.Collection = "email"

	con.Connect()
}

func createEmailDetail(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	if r.Method != "POST" {

		respondWithError(w, http.StatusBadRequest, "Invalid Method")
		return
	}

	var emailData modelData.EmailModel
	fmt.Println("body:", r.Body)
	if err := json.NewDecoder(r.Body).Decode(&emailData); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	ema := emailData.EmailTo

	if len(ema) == 0 || emailData.EmailBody == "" {
		respondWithError(w, http.StatusBadRequest, "Please enter emailTo and emailBody")
		return
	}
	if emailData.EmailSubject == "" {
		respondWithError(w, http.StatusBadRequest, "Please enter subject")
		return
	}
	if emailsend, err := con.Insert(emailData); err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable To Insert Record")
	} else {
		respondWithJson(w, http.StatusAccepted, map[string]string{
			"message": emailsend,
		})
	}

}

func searchEmailData(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid method")
		return
	}

	var cl modelData.EmailSearch

	if err := json.NewDecoder(r.Body).Decode(&cl); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	fmt.Println(cl)
	if searchdocs, err := con.SearchData(cl); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
	} else {
		respondWithJson(w, http.StatusAccepted, searchdocs)
	}
}

func writeToPDF(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "GET" {
		respondWithError(w, http.StatusBadRequest, "Method not allowed")
	}

	id := strings.Split(r.URL.Path, "/")[2]
	fmt.Println("ID:", id)
	if err := con.WriteEmailDataInPDF(id); err != nil {

		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJson(w, http.StatusAccepted, map[string]string{
			"message": "write pdf data successfully",
		})
	}
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func main() {
	http.HandleFunc("/add-emailRecord/", createEmailDetail)
	http.HandleFunc("/search-emailRecord/", searchEmailData)
	http.HandleFunc("/write-to-pdf/", writeToPDF)

	fmt.Println("Excecuted Main Method")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
