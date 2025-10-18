package game

import "log"

// Game représente l'état complet d'une partie de Power 4
// Contient la grille de jeu, le tour actuel et la position du dernier coup joué
type Game struct {
	Grid     [][]string // Grille 6x7 représentant l'état du jeu
	Turn     string     // Joueur actuel : "R" pour Rouge, "Y" pour Jaune
	LastMove struct {   // Position du dernier coup joué pour l'animation
		Row int // Ligne du dernier coup (0-5, 0=en bas)
		Col int // Colonne du dernier coup (0-6)
	}
}

// NewGame initialise et retourne une nouvelle partie de Power 4
// Crée une grille vide 6x7 et définit le premier joueur sur Rouge
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

	// Retourne un nouveau jeu avec grille vide et premier joueur Rouge
	return &Game{
		Grid: grid, // Grille 6x7 vide
		Turn: "R",  // Premier joueur : Rouge
	}
}
<<<<<<< HEAD

=======
>>>>>>> b9bdf19762ad3d1cd1a9bce8d144b03ce7aa5f5c
// PlayMove - VERSION SIMPLIFIÉE et CORRECTE
// Les jetons tombent vers le BAS (ligne 5 = bas, ligne 0 = haut)
func (g *Game) PlayMove(col int) bool {
	if col < 0 || col >= len(g.Grid[0]) {
		return false
	}
<<<<<<< HEAD

=======
	
>>>>>>> b9bdf19762ad3d1cd1a9bce8d144b03ce7aa5f5c
	// CORRECTION : Parcours de BAS EN HAUT pour trouver la première case vide
	// Ligne 5 = bas, Ligne 0 = haut
	for i := len(g.Grid) - 1; i >= 0; i-- {
		if g.Grid[i][col] == "" {
			// Case vide trouvée, placer le jeton
			g.Grid[i][col] = g.Turn
			g.LastMove.Row = i
			g.LastMove.Col = col
			g.Turn = map[string]string{"R": "Y", "Y": "R"}[g.Turn]
			return true
		}
	}
<<<<<<< HEAD

=======
	
>>>>>>> b9bdf19762ad3d1cd1a9bce8d144b03ce7aa5f5c
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
