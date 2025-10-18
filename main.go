package main

import (
	"html/template" // Pour les templates HTML
	"log"           // Pour les logs
	"net/http"      // Pour le serveur HTTP
	"power4/game"   // Notre package game
	"strconv"       // Pour convertir string vers int
)

// currentGame est la partie en cours, stockée en mémoire du serveur
var currentGame = game.NewGame()

// main est la fonction principale qui lance le serveur HTTP
func main() {
	// Enregistrement des routes (URL) et de leurs fonctions de traitement
	http.HandleFunc("/", homeHandler)            // Page d'accueil
	http.HandleFunc("/game", gameHandler)          // Page principale du jeu
	http.HandleFunc("/play", playHandler)        // Traitement des coups (version classique)
	http.HandleFunc("/reset", resetHandler)      // Réinitialisation de la partie
	http.HandleFunc("/clean", cleanHandler)      // Nettoyage du modal
	http.HandleFunc("/api/play", apiPlayHandler) // API pour jouer un coup (AJAX)

	// Servir les fichiers statiques (CSS, JS, images) depuis le dossier "static"
	// StripPrefix enlève "/static/" du chemin pour trouver le fichier
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Démarrage du serveur sur le port 8080
	log.Println("Serveur lancé sur http://localhost:8080")

	// Lancement du serveur, bloquant jusqu'à arrêt
	// ListenAndServe écoute sur le port 8080 et sert les requêtes
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// reverseGrid inverse l'ordre des lignes pour l'affichage
// La grille est stockée ligne 0=en bas, mais affichée ligne 0=en haut
func reverseGrid(grid [][]string) [][]string {
	reversed := make([][]string, len(grid))
	for i := range grid {
		reversed[i] = grid[len(grid)-1-i]
	}
	return reversed
}

// homeHandler gère la page d'accueil
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Si on est déjà sur la page de jeu, rediriger
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
	})

	tmpl = template.Must(tmpl.ParseFiles("templates/index.html"))
	
	winner := ""
	if currentGame.CheckWinner() != "" {
		winner = currentGame.CheckWinner()
	} else if currentGame.IsDraw() {
		winner = "Draw"
	}

	tmpl.Execute(w, struct {
		Grid     [][]string
		Turn     string
		Winner   string
		LastMove struct {
			Row int
			Col int
		}
	}{
		Grid:     reverseGrid(currentGame.Grid),
		Turn:     currentGame.Turn,
		Winner:   winner,
		LastMove: currentGame.LastMove,
	})
}



// apiPlayHandler gère les coups joués via AJAX
// Retourne du JSON au lieu de rediriger vers une page HTML
func apiPlayHandler(w http.ResponseWriter, r *http.Request) {
	// Vérification que la méthode est bien POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return // Arrête le traitement
	}

	// Analyse des données du formulaire
	r.ParseForm()
	colStr := r.FormValue("column")  // Récupère la colonne comme string
	col, err := strconv.Atoi(colStr) // Convertit en int

	// Vérification de la validité de la colonne
	if err != nil || col < 0 || col >= 7 {
		http.Error(w, "Colonne invalide", http.StatusBadRequest)
		return // Colonne doit être entre 0 et 6
	}

	// Vérification que la partie n'est pas déjà terminée
	if currentGame.CheckWinner() != "" || currentGame.IsDraw() {
		http.Error(w, "Partie terminée", http.StatusBadRequest)
		return // On ne peut pas jouer si partie finie
	}

	// Tentative de jouer le coup
	success := currentGame.PlayMove(col)
	if !success {
		http.Error(w, "Colonne pleine", http.StatusBadRequest)
		return // La colonne est pleine
	}

	// Log du coup joué
	log.Println("Coup joué via API dans la colonne :", col)

	// Préparation de la réponse JSON
	w.Header().Set("Content-Type", "application/json")

	// Construction manuelle du JSON pour éviter les imports supplémentaires
	// Format: {"success": true, "lastMove": {"row": X, "col": Y}}
	response := `{"success": true, "lastMove": {"row": ` + strconv.Itoa(currentGame.LastMove.Row) + `, "col": ` + strconv.Itoa(currentGame.LastMove.Col) + }}

	// Envoi de la réponse
	w.Write([]byte(response))
}

// playHandler gère les coups joués via formulaire classique (fallback)
// Redirige vers la page principale après le coup
func playHandler(w http.ResponseWriter, r *http.Request) {
	// Vérification que c'est bien une requête POST
	if r.Method == http.MethodPost {
		// Analyse des données du formulaire
		r.ParseForm()
		colStr := r.FormValue("column")  // Colonne choisie
		col, err := strconv.Atoi(colStr) // Conversion en int

		// Si colonne valide, joue le coup
		if err == nil && col >= 0 && col < 7 {
			currentGame.PlayMove(col) // Joue le coup
			log.Println("Coup joué dans la colonne :", col)
		}
	}

	// Redirection vers la page principale pour voir le résultat
	// StatusSeeOther = 303, force le rechargement
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

// AJOUTEZ cette nouvelle route pour "nettoyer" le modal
func cleanHandler(w http.ResponseWriter, r *http.Request) {
    // Réinitialiser le jeu et rediriger vers l'accueil
    currentGame = game.NewGame()
    http.Redirect(w, r, "/", http.StatusSeeOther)
}