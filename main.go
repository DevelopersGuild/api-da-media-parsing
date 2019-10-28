package main

import (
	"fmt"
	"net/http"
)

func uploadFile(w http.ResponseWriter, r *http.Request){

}

func setupRoutes(){
	http.HandleFunc("/upload", uploadFile)
    http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("Hello World")
	setupRoutes()
}
