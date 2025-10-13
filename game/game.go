package game

type Game struct {
	Grid [][]string
	Turn string
}

func NewGame() *Game {
	grid := make([][]string, 6)
	for i := range grid {
		grid[i] = make([]string, 7)
	}
	return &Game{
		Grid: grid,
		Turn: "R", // Rouge commence
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