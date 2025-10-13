package main

import (
	"log"
	"net/http"
	"html/template"
	"strconv"
	"power4/game"
)

var currentGame = game.NewGame()

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/play", playHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	_ = strconv.Itoa(42) // utilisation temporaire de strconv

	log.Println("Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}
func playHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		colStr := r.FormValue("column")
		col, err := strconv.Atoi(colStr)
		if err != nil {
			log.Println("Coup joué dans la colonne :", col)

		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	
}
