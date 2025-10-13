package main

import (
	"html/template"
	"log"
	"net/http"
	"power4/game"
	"strconv"
)

var currentGame = game.NewGame()

func main() {
	http.HandleFunc("/", handler)           // Affiche la grille
	http.HandleFunc("/play", playHandler)   // Joue un coup
	http.HandleFunc("/reset", resetHandler) // Réinitialise la partie
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// 🔁 Inverse la grille pour affichage du bas vers le haut
func reverseGrid(grid [][]string) [][]string {
	reversed := make([][]string, len(grid))
	for i := range grid {
		reversed[i] = grid[len(grid)-1-i]
	}
	return reversed
}

// 🧠 Affiche la page HTML avec la grille et les boutons
func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Grille actuelle :", currentGame.Grid) // ✅ ici c’est bon

	tmpl := template.New("index.html").Funcs(template.FuncMap{
		"seq": func(start, end int) []int {
			s := make([]int, end-start+1)
			for i := range s {
				s[i] = start + i
			}
			return s
		},
	})
	tmpl = template.Must(tmpl.ParseFiles("templates/index.html"))
	tmpl.Execute(w, struct{ Grid [][]string }{Grid: reverseGrid(currentGame.Grid)})
}

// 🎮 Joue un coup dans la colonne choisie
func playHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		colStr := r.FormValue("column")
		col, err := strconv.Atoi(colStr)
		if err == nil && col >= 0 && col < 7 {
			currentGame.PlayMove(col)
			log.Println("Coup joué dans la colonne :", col)
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// 🔄 Réinitialise la partie
func resetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		currentGame = game.NewGame()
		log.Println("Nouvelle partie lancée")
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
