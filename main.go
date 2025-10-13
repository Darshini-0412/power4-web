go
// main.go
package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var game *Game

func init() {
	game = NewGame()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	data := struct {
		Grid         [6][7]int
		CurrentTurn  int
		GameOver     bool
		Message      string
		ValidColumns []int
	}{
		Grid:         game.Grid,
		CurrentTurn:  game.CurrentTurn,
		GameOver:     game.GameOver,
		Message:      game.Message,
		ValidColumns: game.ValidColumns(),
	}
	tmpl.Execute(w, data)
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	colStr := r.FormValue("col")
	col, err := strconv.Atoi(colStr)
	if err != nil {
		http.Error(w, "Colonne invalide", http.StatusBadRequest)
		return
	}

	err = game.Play(col)
	if err != nil {
		// On garde le message d'erreur simple → on recharge la page normalement
		// (le jeu affichera automatiquement l'état actuel)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		game = NewGame()
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/play", playHandler)
	http.HandleFunc("/reset", resetHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	log.Println(".Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
