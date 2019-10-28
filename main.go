package main

import (
	"context"
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
	// create gcloud client
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		fmt.Println(err)
	}
	bucketName := "api-da-test-bucket"
	enableCors(&w)
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("fileUpload")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	println(fileBytes)
	wc := client.Bucket(bucketName).Object(handler.Filename).NewWriter(ctx)
	if _, err = io.Copy(wc, file); err != nil {
		fmt.Println(err)
	}
	if err := wc.Close(); err != nil {
		fmt.Println(err)
	}

	opts := &storage.SignedURLOptions{
		Expires: time.Date(2025, 12, 22, 0, 0, 0, 0, time.UTC),
		GoogleAccessID: "api-da-housing-test@crack-producer-252518.iam.gserviceaccount.com",
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
