package game

import "log"

// Player représente un joueur avec son nom et sa couleur personnalisée
type Player struct {
	Name   string // Nom du joueur
	Color  string // Couleur personnalisée en hexadécimal (ex: "#FF0000")
	Symbol string // Symbole pour l'identification interne ("R" ou "Y")
}

// Game représente l'état complet d'une partie de Power 4
// Contient la grille de jeu, le tour actuel, les joueurs et la position du dernier coup joué
type Game struct {
	Grid     [][]string // Grille 6x7 représentant l'état du jeu
	Turn     string     // Joueur actuel : "R" pour Rouge, "Y" pour Jaune
	Player1  Player     // Informations du joueur 1
	Player2  Player     // Informations du joueur 2
	LastMove struct {   // Position du dernier coup joué pour l'animation
		Row int // Ligne du dernier coup (0-5, 0=en bas)
		Col int // Colonne du dernier coup (0-6)
	}
}

// NewGame initialise et retourne une nouvelle partie de Power 4
// Crée une grille vide 6x7 et définit les joueurs par défaut
func NewGame() *Game {
	// Log pour le débogage
	log.Println("NewGame() appelé — grille réinitialisée")

	// Création d'une grille 6 lignes x 7 colonnes
	grid := make([][]string, 6)
	for i := range grid {
		// Initialisation de chaque ligne avec 7 colonnes
		grid[i] = make([]string, 7)
		// Remplissage de chaque cellule avec une chaîne vide
		for j := range grid[i] {
			grid[i][j] = "" // Case vide
		}
	}

	// Configuration des joueurs par défaut
	player1 := Player{Name: "Joueur 1", Color: "#FF0000", Symbol: "R"} // Rouge par défaut
	player2 := Player{Name: "Joueur 2", Color: "#FFFF00", Symbol: "Y"} // Jaune par défaut

	// Retourne un nouveau jeu avec grille vide et premier joueur Rouge
	return &Game{
		Grid:    grid, // Grille 6x7 vide
		Turn:    "R",  // Premier joueur : Rouge
		Player1: player1,
		Player2: player2,
	}
}

// SetupPlayers configure les informations des joueurs
// Permet de personnaliser les noms et couleurs des joueurs
func (g *Game) SetupPlayers(player1Name, player1Color, player2Name, player2Color string) {
	// Mise à jour des informations du joueur 1
	g.Player1.Name = player1Name
	g.Player1.Color = player1Color
	g.Player1.Symbol = "R" // Toujours "R" pour le joueur 1

	// Mise à jour des informations du joueur 2
	g.Player2.Name = player2Name
	g.Player2.Color = player2Color
	g.Player2.Symbol = "Y" // Toujours "Y" pour le joueur 2

	log.Printf("Configuration des joueurs: %s (%s) vs %s (%s)",
		g.Player1.Name, g.Player1.Color, g.Player2.Name, g.Player2.Color)
}

// GetCurrentPlayer retourne les informations du joueur actuel
// Utile pour l'affichage dans l'interface utilisateur
func (g *Game) GetCurrentPlayer() Player {
	if g.Turn == "R" {
		return g.Player1
	}
	return g.Player2
}

// GetPlayerBySymbol retourne un joueur basé sur son symbole ("R" ou "Y")
// Permet de récupérer les informations d'un joueur spécifique
func (g *Game) GetPlayerBySymbol(symbol string) Player {
	if symbol == "R" {
		return g.Player1
	}
	return g.Player2
}

// PlayMove - VERSION SIMPLIFIÉE et CORRECTE
// Les jetons tombent vers le BAS (ligne 5 = bas, ligne 0 = haut)
func (g *Game) PlayMove(col int) bool {
	// Vérification de la validité de la colonne
	if col < 0 || col >= len(g.Grid[0]) {
		return false
	}

	// CORRECTION : Parcours de BAS EN HAUT pour trouver la première case vide
	// Ligne 5 = bas, Ligne 0 = haut
	for i := len(g.Grid) - 1; i >= 0; i-- {
		if g.Grid[i][col] == "" {
			// Case vide trouvée, placer le jeton
			g.Grid[i][col] = g.Turn
			g.LastMove.Row = i
			g.LastMove.Col = col

			// Changement de tour
			g.Turn = map[string]string{"R": "Y", "Y": "R"}[g.Turn]
			return true
		}
	}

	return false // Colonne pleine
}

// CheckWinner vérifie s'il y a un gagnant dans la grille
// Retourne "R" si Rouge gagne, "Y" si Jaune gagne, "" sinon
func (g *Game) CheckWinner() string {
	// Dimensions de la grille
	rows := len(g.Grid)    // 6 lignes
	cols := len(g.Grid[0]) // 7 colonnes

	// Directions de vérification pour un alignement de 4 jetons
	// Chaque direction est un vecteur [deltaLigne, deltaColonne]
	directions := [][2]int{
		{0, 1},  // Horizontal → (même ligne, colonne +1)
		{1, 0},  // Vertical ↓ (ligne +1, même colonne)
		{1, 1},  // Diagonale ↘ (ligne +1, colonne +1)
		{-1, 1}, // Diagonale ↗ (ligne -1, colonne +1)
	}

	// Parcours de chaque cellule de la grille
	for r := 0; r < rows; r++ { // Pour chaque ligne
		for c := 0; c < cols; c++ { // Pour chaque colonne
			player := g.Grid[r][c] // Joueur dans la cellule actuelle

			// Ignorer les cases vides
			if player == "" {
				continue // Passe à la cellule suivante
			}

			// Vérification dans les 4 directions
			for _, d := range directions {
				count := 1 // Compte le jeton actuel

				// Vérifie les 3 cases suivantes dans cette direction
				for i := 1; i < 4; i++ {
					// Calcul des coordonnées de la case à vérifier
					nr := r + d[0]*i // Nouvelle ligne
					nc := c + d[1]*i // Nouvelle colonne

					// Vérification des limites de la grille
					if nr < 0 || nr >= rows || nc < 0 || nc >= cols {
						break // Sort de la boucle si hors limites
					}

					// Vérifie si la case contient le même joueur
					if g.Grid[nr][nc] == player {
						count++ // Incrémente le compteur d'alignement
					} else {
						break // Arrête si interruption de l'alignement
					}
				}

				// Si 4 jetons alignés, on a un gagnant
				if count == 4 {
					return player // Retourne "R" ou "Y"
				}
			}
		}
	}

	// Aucun gagnant trouvé
	return "" // Partie continue
}

// IsDraw vérifie si la partie est nulle (grille pleine sans gagnant)
// Retourne true si match nul, false sinon
func (g *Game) IsDraw() bool {
	// Si il y a déjà un gagnant, ce n'est pas un match nul
	if g.CheckWinner() != "" {
		return false // Il y a un gagnant, donc pas de match nul
	}

	// Vérifie si la grille est complètement remplie
	// Parcourt chaque cellule de la grille
	for _, row := range g.Grid { // Pour chaque ligne
		for _, cell := range row { // Pour chaque cellule de la ligne
			if cell == "" {
				// Si une case est vide, la grille n'est pas pleine
				return false // Partie peut continuer
			}
		}
	}

	// Toutes les cases sont remplies et aucun gagnant
	return true // Match nul
}
