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

// 🔁 Inverse les lignes de la grille pour affichage de bas en haut
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
			s := []int{}
			step := 1
			if start > end {
				step = -1
			}
			for i := start; i != end+step; i += step {
				s = append(s, i)
			}
		return s
	},
	})

	tmpl = template.Must(tmpl.ParseFiles("templates/index.html"))
	tmpl.Execute(w, struct {
		Grid   [][]string
		Turn   string
		Winner string
	}{
		Grid:   reverseGrid(currentGame.Grid),
		Turn:   currentGame.Turn,
		Winner: func() string {
			if currentGame.CheckWinner() != "" {
				return currentGame.CheckWinner()
			} else if currentGame.IsDraw() {
				return "Draw"
			}
			return ""
		}(),
	})
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
