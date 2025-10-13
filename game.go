go
// game.go
package main

import (
	"errors"
	"fmt"
)

const (
	ROWS    = 6
	COLS    = 7
	EMPTY   = 0
	PLAYER1 = 1
	PLAYER2 = 2
)

type Game struct {
	Grid        [ROWS][COLS]int
	CurrentTurn int
	GameOver    bool
	Winner      int
	Message     string
}

func NewGame() *Game {
	g := &Game{
		CurrentTurn: PLAYER1,
		GameOver:    false,
		Winner:      EMPTY,
		Message:     "Joueur 1, à toi de jouer !",
	}
	return g
}

func (g *Game) Play(col int) error {
	if g.GameOver {
		return errors.New("la partie est terminée")
	}
	if col < 1 || col > COLS {
		return fmt.Errorf("colonne invalide : %d (doit être entre 1 et %d)", col, COLS)
	}

	// Convertir en index 0-based
	c := col - 1

	// Trouver la première ligne vide (du bas vers le haut)
	for r := ROWS - 1; r >= 0; r-- {
		if g.Grid[r][c] == EMPTY {
			g.Grid[r][c] = g.CurrentTurn
			if g.checkWin(r, c) {
				g.GameOver = true
				g.Winner = g.CurrentTurn
				g.Message = fmt.Sprintf("🎉 Joueur %d gagne !", g.CurrentTurn)
			} else if g.isBoardFull() {
				g.GameOver = true
				g.Message = "🤝 Match nul !"
			} else {
				g.CurrentTurn = 3 - g.CurrentTurn // Alterne 1 ↔ 2
				g.Message = fmt.Sprintf("Joueur %d, à toi de jouer !", g.CurrentTurn)
			}
			return nil
		}
	}
	return errors.New("colonne pleine")
}

func (g *Game) checkWin(row, col int) bool {
	player := g.Grid[row][col]
	directions := [][2]int{
		{0, 1},   // horizontal
		{1, 0},   // vertical
		{1, 1},   // diagonal ↘
		{1, -1},  // diagonal ↙
	}

	for _, d := range directions {
		count := 1
		// Positif
		for i := 1; i < 4; i++ {
			r, c := row+i*d[0], col+i*d[1]
			if r >= 0 && r < ROWS && c >= 0 && c < COLS && g.Grid[r][c] == player {
				count++
			} else {
				break
			}
		}
		// Négatif
		for i := 1; i < 4; i++ {
			r, c := row-i*d[0], col-i*d[1]
			if r >= 0 && r < ROWS && c >= 0 && c < COLS && g.Grid[r][c] == player {
				count++
			} else {
				break
			}
		}
		if count >= 4 {
			return true
		}
	}
	return false
}

func (g *Game) isBoardFull() bool {
	for c := 0; c < COLS; c++ {
		if g.Grid[0][c] == EMPTY {
			return false
		}
	}
	return true
}

func (g *Game) ValidColumns() []int {
	var cols []int
	for c := 0; c < COLS; c++ {
		if g.Grid[0][c] == EMPTY {
			cols = append(cols, c+1)
		}
	}
	return cols
}