package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/play", playHandler)

	http.Handle("/static/", http.StripPrefix("/static/", httpFileServer(http.Dir("static"))))

	log.Println("Serveur lancé sur http://localhost:")
}
