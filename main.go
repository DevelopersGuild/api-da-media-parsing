package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func uploadFile(w http.ResponseWriter, r *http.Request) {	
	enableCors(&w)
	// create gcloud client
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		fmt.Println(err)
	}
	bucketName := "api-da-test-bucket"
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
	objectattrs := storage.ObjectAttrs{
		ContentType: handler.Header.Get("Content-Type"),
		Name:        handler.Filename,
	}
	wc := client.Bucket(bucketName).Object(handler.Filename).NewWriter(ctx)
	wc.ObjectAttrs = objectattrs
	if _, err = io.Copy(wc, file); err != nil {
		fmt.Println(err)
	}
	if err := wc.Close(); err != nil {
		fmt.Println(err)
	}

	// JSON STUFF
	jsonFile, err := os.Open("./auth.json")
	defer jsonFile.Close()
	jsonByteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal([]byte(jsonByteValue), &result)

	if err != nil {
		fmt.Println(err)
	}

	opts := &storage.SignedURLOptions{
		Expires:        time.Date(2025, 12, 22, 0, 0, 0, 0, time.UTC),
		GoogleAccessID: "api-da-housing-test@crack-producer-252518.iam.gserviceaccount.com",
		ContentType:    handler.Header.Get("Content-Type"),
		PrivateKey:     []byte(result["private_key"].(string)),
		Method:         "GET",
	}
	url, err := storage.SignedURL(bucketName, handler.Filename, opts)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Fprintf(w, url)
}

func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8080", nil)
}

func main() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./auth.json")
	fmt.Println("Server running on PORT:8080")
	setupRoutes()
}
