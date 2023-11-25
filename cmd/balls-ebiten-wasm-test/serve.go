package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("../../docs"))
	http.Handle("/", fs)

	log.Println("Listening on http://localhost:8080/index.html")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
