package game

import "log"

type Game struct {
	Grid [][]string
	Turn string
}

func NewGame() *Game {
	log.Println("NewGame() appelé — grille réinitialisée")
	grid := make([][]string, 6)
	for i := range grid {
		grid[i] = make([]string, 7)
	}
	return &Game{
		Grid: grid,
		Turn: "R",
	}
}

func (g *Game) PlayMove(col int) {
	for i := len(g.Grid) - 1; i >= 0; i-- {
		if g.Grid[i][col] == "" {
			g.Grid[i][col] = g.Turn
			g.Turn = map[string]string{"R": "Y", "Y": "R"}[g.Turn]
			break
		}
	}
}

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

func (g *Game) IsDraw() bool {
		// s'il y'a un gagnant, ce n'est pas un match nul
	if g.CheckWinner() != "" {
		return false
	}
	// vérifier si la grille est pleine
	for _, row := range g.Grid {
		for _, cell := range row {
			if cell == "" {
				return false
			}
		}
	}
	return true
}
