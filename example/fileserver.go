package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", http.FileServer(http.Dir(".")))
    http.ListenAndServe(":8080", nil)
}
