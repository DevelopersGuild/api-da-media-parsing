package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func uploadFile(w http.ResponseWriter, r *http.Request){
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
	fmt.Fprintf(w, handler.Filename)
}

func setupRoutes(){
	http.HandleFunc("/upload", uploadFile)
    http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("Server running on PORT:8080")
	setupRoutes()
}
