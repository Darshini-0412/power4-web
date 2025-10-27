package game

import "log"

// Player représente un joueur
type Player struct {
	Name   string
	Color  string
	Symbol string
}

// Game représente l'état complet d'une partie de Power 4
type Game struct {
	Grid     [][]string
	Turn     string
	Player1  Player
	Player2  Player
	LastMove struct {
		Row int
		Col int
	}
}

// NewGame initialise et retourne une nouvelle partie
func NewGame() *Game {
	log.Println("NewGame() appelé — grille réinitialisée")

	grid := make([][]string, 6)
	for i := range grid {
		grid[i] = make([]string, 7)
		for j := range grid[i] {
			grid[i][j] = ""
		}
	}

	player1 := Player{Name: "Joueur 1", Color: "#FF0000", Symbol: "R"}
	player2 := Player{Name: "Joueur 2", Color: "#FFFF00", Symbol: "Y"}

	return &Game{
		Grid:    grid,
		Turn:    "R",
		Player1: player1,
		Player2: player2,
	}
}

// SetupPlayers configure les informations des joueurs
func (g *Game) SetupPlayers(player1Name, player1Color, player2Name, player2Color string) {
	g.Player1.Name = player1Name
	g.Player1.Color = player1Color
	g.Player1.Symbol = "R"

	g.Player2.Name = player2Name
	g.Player2.Color = player2Color
	g.Player2.Symbol = "Y"

	log.Printf("Configuration des joueurs: %s (%s) vs %s (%s)",
		g.Player1.Name, g.Player1.Color, g.Player2.Name, g.Player2.Color)
}

// GetCurrentPlayer retourne le joueur actuel
func (g *Game) GetCurrentPlayer() Player {
	if g.Turn == "R" {
		return g.Player1
	}
	return g.Player2
}

// GetPlayerBySymbol retourne un joueur selon son symbole
func (g *Game) GetPlayerBySymbol(symbol string) Player {
	if symbol == "R" {
		return g.Player1
	}
	return g.Player2
}

// PlayMove joue un coup dans la colonne donnée
func (g *Game) PlayMove(col int) bool {
	if col < 0 || col >= len(g.Grid[0]) {
		return false
	}

	for i := len(g.Grid) - 1; i >= 0; i-- {
		if g.Grid[i][col] == "" {
			g.Grid[i][col] = g.Turn
			g.LastMove.Row = i
			g.LastMove.Col = col
			g.Turn = map[string]string{"R": "Y", "Y": "R"}[g.Turn]
			return true
		}
	}

	return false
}

// CheckWinner vérifie s'il y a un gagnant
func (g *Game) CheckWinner() string {
	rows := len(g.Grid)
	cols := len(g.Grid[0])

	directions := [][2]int{
		{0, 1},  // horizontal
		{1, 0},  // vertical
		{1, 1},  // diagonale ↘
		{-1, 1}, // diagonale ↗
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			player := g.Grid[r][c]
			if player == "" {
				continue
			}
			for _, d := range directions {
				count := 1
				for i := 1; i < 4; i++ {
					nr := r + d[0]*i
					nc := c + d[1]*i
					if nr < 0 || nr >= rows || nc < 0 || nc >= cols {
						break
					}
					if g.Grid[nr][nc] == player {
						count++
					} else {
						break
					}
				}
				if count == 4 {
					return player
				}
			}
		}
	}
	return ""
}

// IsDraw vérifie si la partie est nulle
func (g *Game) IsDraw() bool {
	if g.CheckWinner() != "" {
		return false
	}
	for _, row := range g.Grid {
		for _, cell := range row {
			if cell == "" {
				return false
			}
		}
	}
	return true
}
