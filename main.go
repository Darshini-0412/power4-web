package main

import (
	"html/template" // Pour les templates HTML
	"log"           // Pour les logs
	"net/http"      // Pour le serveur HTTP
	"power4/game"   // Notre package game
	"strconv"       // Pour convertir string vers int
	"strings"
)

// currentGame est la partie en cours, stockée en mémoire du serveur
var currentGame = game.NewGame()

// main est la fonction principale qui lance le serveur HTTP
func main() {
	// Enregistrement des routes (URL) et de leurs fonctions de traitement
	http.HandleFunc("/", homeHandler)            // Page d'accueil
	http.HandleFunc("/game", gameHandler)        // Page principale du jeu
	http.HandleFunc("/play", playHandler)        // Traitement des coups (version classique)
	http.HandleFunc("/reset", resetHandler)      // Réinitialisation de la partie
	http.HandleFunc("/clean", cleanHandler)      // Nettoyage du modal
	http.HandleFunc("/api/play", apiPlayHandler) // API pour jouer un coup (AJAX)
	http.HandleFunc("/api/setup", setupHandler)  // API pour configurer les joueurs

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

	// Récupération des informations des joueurs pour le template
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

	// Analyse des données du formulaire
	if err := r.ParseForm(); err != nil {
		log.Printf("Erreur ParseForm: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "message": "Erreur d'analyse du formulaire"}`))
		return
	}

	// DEBUG: Afficher tous les paramètres reçus
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

	// Si les couleurs sont vides, essayer de les récupérer différemment
	if player1Color == "" {
		// Essayer de récupérer depuis les valeurs du formulaire
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

	// Si toujours vide, utiliser les valeurs par défaut
	if player1Color == "" {
		player1Color = "#FF0000"
		log.Printf("Couleur joueur 1 vide, utilisation de la valeur par défaut: %s", player1Color)
	}

	if player2Color == "" {
		player2Color = "#FFFF00"
		log.Printf("Couleur joueur 2 vide, utilisation de la valeur par défaut: %s", player2Color)
	}

	// Normalisation des couleurs
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

	// Réponse JSON de succès
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true, "message": "Configuration sauvegardée"}`))

	log.Printf("Joueurs configurés avec succès: %s (%s) vs %s (%s)",
		player1Name, player1Color, player2Name, player2Color)
}

// isValidColor vérifie si une couleur est valide (format hexadécimal)
func isValidColor(color string) bool {
	if color == "" {
		return false
	}

	// Accepter les formats: #RGB, #RRGGBB
	if len(color) != 7 && len(color) != 4 {
		return false
	}

	if color[0] != '#' {
		return false
	}

	// Vérifier que chaque caractère est un chiffre hexadécimal
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
	response := `{"success": true, "lastMove": {"row": ` + strconv.Itoa(currentGame.LastMove.Row) + `, "col": ` + strconv.Itoa(currentGame.LastMove.Col) + `}}`

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
