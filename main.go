package main

import (
	"html/template"
	"log"
	"net/http"
	"power4/game"
	"strconv"
	"strings"
)

var currentGame = game.NewGame()

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/game", gameHandler)
	http.HandleFunc("/play", playHandler)
	http.HandleFunc("/reset", resetHandler)
	http.HandleFunc("/api/play", apiPlayHandler)
	http.HandleFunc("/api/setup", setupHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Démarrage du serveur sur le port 8080
	log.Println("Serveur lancé sur http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func reverseGrid(grid [][]string) [][]string {
	reversed := make([][]string, len(grid))
	for i := range grid {
		reversed[i] = grid[len(grid)-1-i]
	}
	return reversed
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, nil)
}

// gameHandler gère la page de jeu (anciennement handler)
func gameHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Grille actuelle :", currentGame.Grid)

	tmpl := template.New("index.html").Funcs(template.FuncMap{
		"seq": func(start, end int) []int {
			var s []int
			if start <= end {
				for i := start; i <= end; i++ {
					s = append(s, i)
				}
			} else {
				for i := start; i >= end; i-- {
					s = append(s, i)
				}
			}
			return s
		},
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
	})

	tmpl = template.Must(tmpl.ParseFiles("templates/index.html"))

	winner := ""
	if currentGame.CheckWinner() != "" {
		winner = currentGame.CheckWinner()
	} else if currentGame.IsDraw() {
		winner = "Draw"
	}

	player1 := currentGame.Player1
	player2 := currentGame.Player2

	tmpl.Execute(w, struct {
		Grid     [][]string
		Turn     string
		Winner   string
		LastMove struct {
			Row int
			Col int
		}
		Player1 game.Player
		Player2 game.Player
	}{
		Grid:     reverseGrid(currentGame.Grid),
		Turn:     currentGame.Turn,
		Winner:   winner,
		LastMove: currentGame.LastMove,
		Player1:  player1,
		Player2:  player2,
	})
}

// setupHandler gère la configuration des joueurs
func setupHandler(w http.ResponseWriter, r *http.Request) {
	// Vérification que la méthode est bien POST
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"success": false, "message": "Method not allowed"}`))
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("Erreur ParseForm: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "message": "Erreur d'analyse du formulaire"}`))
		return
	}

	log.Printf("=== DEBUG FORM DATA ===")
	for key, values := range r.Form {
		log.Printf("  %s: %v", key, values)
	}
	log.Printf("=== FIN DEBUG ===")

	// Récupération des données
	player1Name := r.FormValue("player1Name")
	player1Color := r.FormValue("player1Color")

	player2Name := r.FormValue("player2Name")
	player2Color := r.FormValue("player2Color")

	log.Printf("DEBUG - player1Name: '%s'", player1Name)
	log.Printf("DEBUG - player1Color: '%s'", player1Color)
	log.Printf("DEBUG - player2Name: '%s'", player2Name)
	log.Printf("DEBUG - player2Color: '%s'", player2Color)

	if player1Color == "" {
		if colors, ok := r.Form["player1Color"]; ok && len(colors) > 0 {
			player1Color = colors[0]
			log.Printf("DEBUG - player1Color récupéré depuis Form[]: '%s'", player1Color)
		}
	}

	if player2Color == "" {
		if colors, ok := r.Form["player2Color"]; ok && len(colors) > 0 {
			player2Color = colors[0]
			log.Printf("DEBUG - player2Color récupéré depuis Form[]: '%s'", player2Color)
		}
	}

	if player1Color == "" {
		player1Color = "#FF0000"
		log.Printf("Couleur joueur 1 vide, utilisation de la valeur par défaut: %s", player1Color)
	}

	if player2Color == "" {
		player2Color = "#FFFF00"
		log.Printf("Couleur joueur 2 vide, utilisation de la valeur par défaut: %s", player2Color)
	}

	player1Color = strings.ToUpper(player1Color)
	player2Color = strings.ToUpper(player2Color)

	log.Printf("Configuration finale - Joueur1: %s (%s), Joueur2: %s (%s)",
		player1Name, player1Color, player2Name, player2Color)

	// Validation des couleurs
	if !isValidColor(player1Color) {
		log.Printf("Couleur joueur 1 invalide: %s", player1Color)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "message": "Couleur du joueur 1 invalide: ` + player1Color + `"}`))
		return
	}

	if !isValidColor(player2Color) {
		log.Printf("Couleur joueur 2 invalide: %s", player2Color)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "message": "Couleur du joueur 2 invalide: ` + player2Color + `"}`))
		return
	}

	// Validation des noms
	if strings.TrimSpace(player1Name) == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "message": "Nom du joueur 1 invalide"}`))
		return
	}

	if strings.TrimSpace(player2Name) == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "message": "Nom du joueur 2 invalide"}`))
		return
	}

	// Validation que les couleurs sont différentes
	if player1Color == player2Color {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "message": "Les couleurs doivent être différentes"}`))
		return
	}

	// Configuration des joueurs dans le jeu
	currentGame.SetupPlayers(player1Name, player1Color, player2Name, player2Color)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true, "message": "Configuration sauvegardée"}`))

	log.Printf("Joueurs configurés avec succès: %s (%s) vs %s (%s)",
		player1Name, player1Color, player2Name, player2Color)
}

func isValidColor(color string) bool {
	if color == "" {
		return false
	}

	if len(color) != 7 && len(color) != 4 {
		return false
	}

	if color[0] != '#' {
		return false
	}

	colorPart := color[1:]
	for _, char := range colorPart {
		if !((char >= '0' && char <= '9') ||
			(char >= 'a' && char <= 'f') ||
			(char >= 'A' && char <= 'F')) {
			return false
		}
	}

	return true
}

func apiPlayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	colStr := r.FormValue("column")
	col, err := strconv.Atoi(colStr)

	if err != nil || col < 0 || col >= 7 {
		http.Error(w, "Colonne invalide", http.StatusBadRequest)
		return
	}

	if currentGame.CheckWinner() != "" || currentGame.IsDraw() {
		http.Error(w, "Partie terminée", http.StatusBadRequest)
		return
	}

	// Tentative de jouer le coup
	success := currentGame.PlayMove(col)
	if !success {
		http.Error(w, "Colonne pleine", http.StatusBadRequest)
		return
	}

	log.Println("Coup joué via API dans la colonne :", col)

	w.Header().Set("Content-Type", "application/json")

	response := `{"success": true, "lastMove": {"row": ` + strconv.Itoa(currentGame.LastMove.Row) + `, "col": ` + strconv.Itoa(currentGame.LastMove.Col) + `}}`

	w.Write([]byte(response))
}

// playHandler gère les coups joués
func playHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		r.ParseForm()
		colStr := r.FormValue("column")
		col, err := strconv.Atoi(colStr)

		if err == nil && col >= 0 && col < 7 {
			currentGame.PlayMove(col) // Joue le coup
			log.Println("Coup joué dans la colonne :", col)
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// resetHandler réinitialise la partie en cours
func resetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		currentGame = game.NewGame()
		log.Println("Nouvelle partie lancée")
	}
	http.Redirect(w, r, "/game", http.StatusSeeOther)
}

func cleanHandler(w http.ResponseWriter, r *http.Request) {
	currentGame = game.NewGame()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
