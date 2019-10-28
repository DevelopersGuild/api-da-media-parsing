package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	pictureURLChannel := make(chan string)
	go func() {
		// create gcloud client
		ctx := context.Background()
		client, err := storage.NewClient(ctx)
		bucketName := "api-da-test-bucket"

		if err != nil {
			fmt.Println(err)
		}

		// parse file from post requests.
		r.ParseMultipartForm(10 << 20)
		file, handler, err := r.FormFile("fileUpload")

		if err != nil {
			fmt.Println(err)
		}

		defer file.Close()
		// fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}
		// name uuid generation
		fileExtension := strings.Split(handler.Filename, ".")[1]
		uniqueID, err := uuid.NewRandom()
		if err != nil {
			fmt.Println(err)
		}
		uniqueIDString := uniqueID.String()
		// Writer Object Change
		objectattrs := storage.ObjectAttrs{
			ContentType: handler.Header.Get("Content-Type"),
			Name:        uniqueIDString + "." + fileExtension,
		}
		// handler.Filename
		wc := client.Bucket(bucketName).Object(uniqueIDString + "." + fileExtension).NewWriter(ctx)
		wc.ObjectAttrs = objectattrs
		if _, err = io.Copy(wc, file); err != nil {
			fmt.Println(err)
		}
		if err := wc.Close(); err != nil {
			fmt.Println(err)
		}

		// JSON FOR URL
		jsonFile, err := os.Open("./auth.json")
		defer jsonFile.Close()
		jsonByteValue, _ := ioutil.ReadAll(jsonFile)
		var result map[string]interface{}
		json.Unmarshal([]byte(jsonByteValue), &result)

		if err != nil {
			fmt.Println(err)
		}
		pictureURLChannel <- "https://storage.googleapis.com/" + bucketName + "/" + uniqueIDString + "." + fileExtension
	}()
	fmt.Fprintf(w, <-pictureURLChannel)
}

func deleteFile(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	params := mux.Vars(r)
	completionChannel := make(chan string)
	go func() {
		ctx := context.Background()
		client, err := storage.NewClient(ctx)
		bucketName := "api-da-test-bucket"
		if err != nil {
			fmt.Println(err)
		}

		o := client.Bucket(bucketName).Object(params["name"])
		if err := o.Delete(ctx); err != nil {
			fmt.Println(err)
		}
		completionChannel <- params["name"] + " deleted!"
	}()
	fmt.Fprintf(w, <-completionChannel)
}

func setupRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/upload", uploadFile).Methods("POST")
	r.HandleFunc("/delete", deleteFile).Methods("DELETE")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func main() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./auth.json")
	fmt.Println("Server running on PORT:8080")
	setupRoutes()
}
