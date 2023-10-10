package main

import (
	"log"
	"net/http"
)

func main() {
	_ = Run()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
