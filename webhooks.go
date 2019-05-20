package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	router := httprouter.New()
	router.GET("/", Webhook)
	log.Printf("Listen %s:%s", "0.0.0.0", "8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", loggingHandler(router)))
}
