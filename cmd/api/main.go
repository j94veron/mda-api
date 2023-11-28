package main

import (
	"go-api/logger"
	"log"
	"net/http"
)

func main() {
	logger.InitLogger()
	_ = Run()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
